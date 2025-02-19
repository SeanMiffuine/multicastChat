# **Multicast RPC Chat with Vector Clocks**  

A distributed chat application implemented in Go, designed to demonstrate **multicast RPC communication, vector clocks for causal ordering, and a decentralized node structure**. Each chat instance operates as an independent node, forming a dynamically structured distributed system.  

This project showcases fundamental **distributed systems concepts** such as **logical timekeeping, state consistency, and fault-tolerant message propagation**. It also includes a custom-built **RPC and multicast library**, enabling reliable inter-node communication.  

## **Key Features**  

- **Multicast RPC Communication** – Messages are exchanged through a custom-built RPC and multicast framework, ensuring efficient and scalable message propagation.  
- **Vector Clocks for Causal Consistency** – Each message is timestamped using vector clocks, preserving causal relationships and ensuring a coherent chat history across nodes.  
- **Decentralized Node Structure** – Each chat instance acts as a client-node within a dynamic, self-maintaining node set, eliminating the need for a centralized server.  
- **Fault Tolerance & Scalability** – The system is designed to handle node failures gracefully while supporting multiple concurrent chat instances.  

## **How It Works**  

### 1. Bootstrapping the Network  
- A new chat instance (node) starts up and registers itself within a **distributed node set**.  
- Nodes communicate via **RPC calls**, enabling seamless peer-to-peer messaging.  

### 2. Message Exchange via Multicast  
- Messages are propagated using a custom **multicast library**, ensuring all active nodes receive updates efficiently.  
- Vector clocks are attached to messages, maintaining causal order and avoiding inconsistencies.  

### 3. Handling Node Joins & Failures  
- New nodes are dynamically integrated into the network, updating the distributed state.  
- Failures and removal of nodes are to be implemented in future improvements

## **Installation & Usage**  

```sh
git clone https://github.com/YOUR_USERNAME/multicast-rpc-chat.git  
cd multicast-rpc-chat  
go run src/lib/nodeset/cmd/nodesetd/nodesetd.go   # for server node
go run src/cmd/chat/chat.go                       # for chat instance
```

- Start multiple instances in separate terminals to simulate a distributed environment.  
- Send messages and observe **causally ordered** message delivery across nodes.  

## **Why This Matters in Distributed Systems**  

- **Causal Ordering via Vector Clocks** – Ensures that messages are processed in the correct order across nodes, solving concurrency issues.  
- **Multicast Communication** – Demonstrates efficient peer-to-peer messaging without relying on a centralized broker.  
- **Decentralized Node Structure** – Resilient architecture that adapts to network changes dynamically.  
- **Custom RPC & Multicast Implementation** – Showcases a robust, lightweight approach to inter-process communication in distributed environments.  

## **Potential Use Cases**  

- **Distributed Chat Systems** – Serverless, scalable messaging networks.  
- **Event-Driven Architectures** – Applications requiring **order-preserving** event propagation.  
- **Consensus & Coordination** – Building blocks for more complex distributed protocols.  
