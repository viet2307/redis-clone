# redis-clone

## Overview
- A Redis inspired key-value store supporting multi thread capability

### Why re-inventing the wheel?
- Review my own knowledge on the concepts of:
    - Redis concept (RESP, storage data structures, methods,...)
    - Basic knowledge on OS (virtualization, concurrency and persistence)
    - TCP server from scratch
    - I/O models and concurrency related systemcalls (epoll, kqueue...)

#### TODO
- [x] Multi-threaded TCP server using goroutine
    - [ ] I/O multiplexing model with threadpool
- [ ] RESP protocol parser
- [ ] PING, GET, SET implementation
- [ ] Sorted-set implementation

## REDIS
- A in-memory key-value store, stores data in RAM, use AOF (append-only file) command log for persistence
- Use RESP protocol
    - Built on top of TCP
    - Simple to parse
- Redis uses 1 thread **to process user's requests** but do have background thread for other tasks such as logging
- Redis **write to memory first** then append to the AOF log. Trades persistence for speed, data can restore from AOF command log but there can be data loss

## OS
- [Basic OS knowledges](https://pages.cs.wisc.edu/~remzi/OSTEP/#:~:text=free%20online%20form%20of%20the%20book):
    - **Is a process itself**, sleeps in the background and constantly wakes up via interrupt to do tasks
    - **Virtualization** of resources:
        - Helps multiple program **share** physical resources like RAM, CPU,...
        - Abstraction of program via process:
            - Each process is an abstraction of a program, each process is a set of memory addresses wrapped into a sandbox for isolation
        - Abstraction of memory via address space:
            - **Virtual view** of the memory address, isolated for each process, mapped to physical address by OS itself
            - Byte level mapping via page/framing
        - I/O resource list:
            - A data structure that manages **virtual view** of resources it uses like files, sockets, pipes, I/O streams...
        - A process can create another process (**threading**) by invoking system call and tell OS to spawn another process within itself
    - **Concurrency**:
        - Provide concurrency via **context switching** (Tells CPU thread what to run next)
        - Provide **semaphore** for synchronization
            - Semaphore is a kernel resource
            - semaphore's operations are atomic
        - Context switching are broadly seperated into 2 models: **proportional share** model and **Multi-level Feedback Queue** (MLFQ)
        - **Linux** uses proportional share model, processes each have a proportion of execution time.
            - Implemented using red-black tree, each *node* is a task/process, left most nodes is chosen to run (smallest `vruntime`)
            - vruntime changes based on delta time on context switch and `nice` values, which provide tasks priority capability
    - **Persistence**:
        - Handle the mapping between *virtual address/view* of processes to *physical address* of resources on disk, RAM,...
        - Provide data integrity and protection to files, disk

## Application Connection:
- This project will use client-server application model, there are some choices:

![Thread per connection model](assets\client-server.png)

    - Thread per connection model:
        - Pros: simple to implement, leverages multi-core processor, handle blocking I/O
        - Cons: high memory and cpu overhead, risk of race condition

![Thread pool](assets\thread-pool.png)
    
    - Thread pool:
        - Pros: avoid overload harware
        - Cons: hard to configure/choose pool size, high overhead for very short task, more complex and risk of race condition

![Event driven](assets\event-driven.png)

    - Event driven:
        - Pros: scalable for I/O bound app, efficient resource usage (no context switch overhead), reduce race condition
        - Cons: complex, CPU-Bound operations block everything