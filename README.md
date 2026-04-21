# Minimal Container Runtime in Go

This project is a small attempt to understand how containers actually work under the hood.

Instead of relying on Docker as a black box, I built a minimal container runtime from scratch using Go. The goal was not to replicate Docker completely, but to understand the core mechanisms that make containerization possible — process isolation, filesystem isolation, and resource control.



## What this project does

The runtime allows executing a process inside an isolated environment by combining multiple Linux kernel features:

* Process isolation using PID namespaces
* Mount isolation using mount namespaces
* Filesystem isolation using chroot
* Process visibility through a mounted `/proc` filesystem
* Resource control using cgroups (CPU and memory limits)

The result is a minimal but functional container-like environment where processes run independently from the host system.



## How it works

The execution flow is intentionally simple:

* The program starts in "run" mode
* It re-executes itself using `/proc/self/exe` to create a clean process context
* A child process is created with new namespaces
* The process is attached to a cgroup for resource limits
* The root filesystem is changed using `chroot`
* `/proc` is mounted to provide process information
* The target command is executed inside this isolated environment

Each step mirrors a core concept used in real container runtimes.



## Running the project

Build the binary:

```bash
go build -o mydocker
```

Run a command inside the container:

```bash
sudo ./mydocker run /bin/ls /
```

The output should show a minimal filesystem, confirming that the process is running in an isolated root environment.



## Project structure

```
mini-docker/
├── main.go
├── rootfs/
└── README.md
```

The `rootfs` directory acts as the container’s filesystem and must contain the required binaries and their dependencies.



## Limitations

This project focuses on core concepts, so several features are intentionally left out:

* No proper terminal (PTY) support, which affects interactive shells
* No networking support
* No image layering or package management
* Minimal root filesystem created manually

These limitations are a direct result of keeping the implementation simple and focused on learning.



## What I learned
Building this project made a few things very clear:

* Containers are not a single feature, but a combination of multiple Linux primitives
* The filesystem inside a container defines what the process can do
* Isolation requires careful coordination between namespaces, mounts, and process execution
* Many things that feel “automatic” in Docker are actually explicit steps behind the scenes



## Possible improvements

There are several directions this project could be extended:

* Adding networking using virtual interfaces and bridges
* Supporting a more complete root filesystem
* Introducing basic image handling
* Improving terminal support with proper PTY handling
* Adding CLI flags for resource limits



## Summary

This project is not meant to compete with Docker, but to understand it.

It demonstrates how containerization works at a fundamental level by building the core pieces manually. More importantly, it serves as a practical exploration of Linux systems programming and process isolation.



## Resume description

Developed a minimal container runtime in Go implementing Linux namespaces, chroot-based filesystem isolation, and cgroups for resource management. Designed a re-exec based process model and explored system-level constraints such as lack of PTY and networking support.

## Demo (Screen Recorded Video Link)
https://drive.google.com/file/d/1NHiXqpNJzodqo4_Q57JqrOBglV-YK1Xj/view?t=9.394

## Author

Atharv Krishna
