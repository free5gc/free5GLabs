# Lab 6: Git & GitHub

## Introduction

In Lab 6, you need to learn how to use Git and GitHub and utilize them to do a simple task!

This lab does not contain any information about Git!



## The goals of this lab

- Create/Fork a Repository from GitHub and pull it to your local computer. 
- Fork a Repository -> Create Branch -> Create Pull Request to the origin. 
- Follow the Branch Naming Policy and Conventional Commit Message. 
- <font color="red">Knows how to contribute to the free5GC project</font>



## What is Git/GitHub?

> **ChatGTP**
>
> < What is Git? >
>
> >Git is a distributed version control system used to track changes in source code during software development. It allows multiple developers to work on a project simultaneously without interfering with each other’s changes.

- [Git vs. GitHub: What's the difference?](https://www.youtube.com/watch?v=wpISo9TNjfU) -- from IBM
- [Git 的故事：這一次沒這麼好玩](https://blog.brachiosoft.com/posts/git/) -- from brachiosoft

There are plenty of resources explaining what Git, GitHub, and version control are, so we skip the explanation here. 


#### Key Points

- Repository
- Branch
- Commit



## Branch Naming Policy & Conventional Commit Message

> [!CAUTION]
>
> ⚠️ Do not use <font color="red">**main branch**</font> to develop! ⚠️

The description of branch names/commit messages makes it easy to understand the purpose of the branch/commit at a glance. 

For more details, please see [Here](https://hackmd.io/@CTFang/H1TWDLz1A).

#### Branch Name Example

- Lower case only
- Kebab Case

```
feature/oauth-support
refactor/oauth-token-ctx
fix/pdu-release
docs/chf-design
```

#### Conventional Commit Message Example

```
test: add unit tests for OAuth2
refactor: consumer, server, processor
docs: update API documentation
fix: resolve PDU release panic
```



## Exercise

- Fork https://github.com/andy89923/nf-example & Trace code 
    - It's a sample of free5GC NFs 
    - Same structure but no **consumer** 
    - Not a real NF, for sure! 😆
-  Follow the NF structure, create another API service
    - Must include 1 GET and 1 POST method at least
- Create a Pull Request back to the origin repository
- **===== Pending for Pull Request Review(s) =====**
- Revise, Rebase, and Solve Conflicts if have any
- If the Pull Request has been merged, you have finished this Lab! 🎉


