"use client";
import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import axios from "axios";
import { Ticket } from "@/gen/ticket_pb";
import TicketModal from "./[id]/modal";

const PAGE_SIZE = 10;

export default function TicketsPage() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(0);
  const [total, setTotal] = useState(0);

  const searchParams = useSearchParams();
  const selectedId = searchParams.get("id");
  const router = useRouter();

  const totalPages = Math.ceil(total / PAGE_SIZE);

  useEffect(() => {
    axios
      .get<{ tickets: Ticket[]; total: number }>(
        `http://localhost:8081/tickets?limit=${PAGE_SIZE}&offset=${page * PAGE_SIZE}`
      )
      .then((res) => {
        setTickets(res.data.tickets);
        setTotal(res.data.total);
      });
  }, [page]);

  const filtered = tickets.filter(
    (t) =>
      t.customerName.toLowerCase().includes(search.toLowerCase()) ||
      t.email.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="bg-gray-100 text-gray-800 p-6 min-h-screen">
      <div className="max-w-5xl mx-auto relative">
        <h1 className="text-3xl font-bold mb-6">
          🎫 Onboarding Ticket Dashboard
        </h1>

        <input
          type="text"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search by name or email"
          className="w-full px-4 py-2 mb-4 border border-gray-300 rounded shadow-sm bg-white"
        />

        <div className="bg-white shadow rounded overflow-auto max-h-[500px]">
          <table className="w-full text-left border-collapse">
            <thead className="bg-gray-50 sticky top-0">
              <tr>
                <th className="p-3">Customer</th>
                <th className="p-3">Email</th>
                <th className="p-3">Status</th>
                <th className="p-3">Created</th>
                <th className="p-3">Actions</th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((ticket) => (
                <tr key={ticket.id} className="border-t border-gray-300 hover:bg-gray-50">
                  <td className="p-3">{ticket.customerName}</td>
                  <td className="p-3">{ticket.email}</td>
                  <td className="p-3">
                    <span
                      className={`inline-block px-2 py-1 text-sm rounded ${
                        ticket.status === "done"
                          ? "bg-green-200 text-green-800"
                          : ticket.status === "pending"
                          ? "bg-yellow-200 text-yellow-800"
                          : "bg-gray-200 text-gray-800"
                      }`}
                    >
                      {ticket.status}
                    </span>
                  </td>
                  <td className="p-3">{ticket.createdAt.split("T")[0]}</td>
                  <td className="p-3">
                    <button
                      className="text-blue-500 hover:underline"
                      onClick={() => router.push(`/?id=${ticket.id}`)}
                    >
                      Edit
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Pagination controls */}
        <div className="flex justify-between items-center mt-4 text-sm text-gray-700">
          <button
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            disabled={page === 0}
            className="px-4 py-2 bg-gray-200 rounded disabled:opacity-50"
          >
            ← Previous
          </button>

          <span>
            Page {page + 1} of {totalPages || 1}
          </span>

          <button
            onClick={() => setPage((p) => p + 1)}
            disabled={(page + 1) * PAGE_SIZE >= total}
            className="px-4 py-2 bg-gray-200 rounded disabled:opacity-50"
          >
            Next →
          </button>
        </div>

        {selectedId && (
          <TicketModal
            id={Number(selectedId)}
            onClose={() => router.push("/")}
            onUpdated={(updated) => {
              setTickets((tickets) =>
                tickets.map((t) => (t.id === updated.id ? updated : t))
              );
            }}
          />
        )}
      </div>
    </div>
  );
}
