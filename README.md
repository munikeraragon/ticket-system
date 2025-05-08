# ticket-system


### How to run

1. Run ```docker compose build && docker compose up```
2. Go to localhost:3000

###  How to test

Run ```chmod +x ./test.sh && ./test.sh```


### 👋 About Me

I’m a senior engineer who’s worked across several teams leading the design and delivery of scalable, cloud-native applications. I have a strong background in distributed sytems, cloud infrastructure, and security for both containerized and serverless architectures.

### Technical choices

#### Backend

My first technical choice for this project was choosing between serverless and containerized (Kubernetes) architecture. While serverless works well for event-driven workloads and would be suitable for a ticketing system, Halen is a more complex application with high-throughput services like ridesharing, delivery, and logistics. These typically require long lived services, tight network control, and advanced orchestration. All of which are better suited to Kubernetes.

My second choice was language and tooling. I selected Go becuase of its support for concurrency, performance, and production readiness. This aligns with Halen's scale and complexity, and reflects common choices from companies like Uber and Lyft

My third choice was around the API API Protocol. I chose gRPC with gRPC-Gateway to expose both typed protobuf-based APIs and a RESTful interface. gRPC allowed me to:

- Define strong types across backend and frontend
- Auto-generate code from .proto files
- Produce OpenAPI specs for tooling and validation


Lastly, I consider storage. Since this is a transactional system with structured data (tickets), I went with PostgreSQL.


#### Fronent
The frontend was a straightforward choice given my experience with Nextjs + React and its ability to handle:
- Fast development with file-based routing
- API integration via REST
- Scalable UI component using Tailwind CSS

### Next Steps

The next step would be to productionalize and secure the application in a cloud environment.

To start, I’d define Helm charts to standardize the deployment of the Go microservices. Then I’d set up a local Kubernetes cluster to test the services orchestration and networking.

For internal service communication and observability, I’d configure Istio as the service mesh.

Security wise, I’d tighten up the CORS configuration to allow only trusted frontend origins, and set up TLS termination either at the ingress layer or using a managed certificate (via AWS ACM or cert-manager). I’d also implement user authentication using Cognito to protect the dashboard and any sensitive endpoints.

I’d use Terraform to provision and manage the AWS infrastructure (dev, staging, prod), including the VPC, EKS cluster, Cognito for auth, and CloudWatch for centralized logging and metrics.

Finally, I’d use Argo CD to manage deployments across environments. I’d pair that with GitHub Actions to handle builds and image pushes, while letting Argo CD handle rollouts and sync.