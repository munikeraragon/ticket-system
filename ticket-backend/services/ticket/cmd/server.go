package ticket

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "ticket/gen/go"

	_ "github.com/lib/pq"
)

type TicketServer struct {
	pb.UnimplementedTicketServiceServer
	db *sql.DB
}

func NewTicketServer(db *sql.DB) pb.TicketServiceServer {
	return &TicketServer{db: db}
}

func StartGRPCServer(s pb.TicketServiceServer, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTicketServiceServer(grpcServer, s)
	return grpcServer.Serve(lis)
}

func StartHTTPServer(s pb.TicketServiceServer, addr string) error {
	ctx := context.Background()
	mux := runtime.NewServeMux()
	if err := pb.RegisterTicketServiceHandlerServer(ctx, mux, s); err != nil {
		return err
	}
	return http.ListenAndServe(addr, allowCORS(mux))
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow everything during dev — tighten in production
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func (s TicketServer) GetTicket(ctx context.Context, req *pb.GetTicketRequest) (*pb.Ticket, error) {
	log.Printf("GetTicket called with id=%d", req.Id)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, `
		SELECT id, customer_name, email, created_at, status, notes 
		FROM tickets WHERE id = $1`, req.Id)

	var t pb.Ticket
	if err := row.Scan(&t.Id, &t.CustomerName, &t.Email, &t.CreatedAt, &t.Status, &t.Notes); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetTicket: no ticket found for id=%d", req.Id)
			return nil, status.Error(codes.NotFound, "ticket not found")
		}
		log.Printf("GetTicket query error: %v", err)
		return nil, err
	}

	return &t, nil
}

func (s *TicketServer) ListTickets(ctx context.Context, req *pb.ListTicketsRequest) (*pb.ListTicketsResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, customer_name, email, created_at, status, notes 
		FROM tickets
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*pb.Ticket
	for rows.Next() {
		var t pb.Ticket
		if err := rows.Scan(&t.Id, &t.CustomerName, &t.Email, &t.CreatedAt, &t.Status, &t.Notes); err != nil {
			return nil, err
		}
		tickets = append(tickets, &t)
	}

	var total int32
	err = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tickets`).Scan(&total)
	if err != nil {
		return nil, err
	}

	return &pb.ListTicketsResponse{
		Tickets: tickets,
		Total:   total,
	}, nil
}

func (s *TicketServer) UpdateTicket(ctx context.Context, req *pb.UpdateTicketRequest) (*pb.UpdateTicketResponse, error) {
	allowedStatuses := map[string]bool{
		"open":    true,
		"pending": true,
		"done":    true,
	}
	if !allowedStatuses[req.Status] {
		log.Printf("UpdateTicket: invalid status value: %s", req.Status)
		return nil, status.Error(codes.InvalidArgument, "invalid status")
	}

	res, err := s.db.ExecContext(ctx, `
		UPDATE tickets
		SET status = $1, notes = $2
		WHERE id = $3
	`, req.Status, req.Notes, req.Id)

	if err != nil {
		log.Printf("UpdateTicket update error: %v", err)
		return nil, err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("UpdateTicket: no ticket found with id=%d", req.Id)
		return nil, status.Error(codes.NotFound, "ticket not found")
	}

	row := s.db.QueryRowContext(ctx, `
		SELECT id, customer_name, email, created_at, status, notes
		FROM tickets WHERE id = $1
	`, req.Id)

	var t pb.Ticket
	if err := row.Scan(&t.Id, &t.CustomerName, &t.Email, &t.CreatedAt, &t.Status, &t.Notes); err != nil {
		return nil, err
	}

	return &pb.UpdateTicketResponse{UpdatedTicket: &t}, nil
}
