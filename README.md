# worker-pool

Worker-pool is a demo project that show several important techniques

Briefly, worker-pool is a configurable pool of workers.
Worker is something that can run a task. 
And a task is a set of instructions. 

## key techniques implemented in the project
Tasks are created by their configuration, and it's up to Factory to create a Task.
It makes it easy to Mock the interface in tests and replace DefaultFactory with a custom one
that might be configured to build whatever we want. That is extremely useful for testing.

Everywhere in the workflow the Task interface is used instead of particular Task implementations.
That simplifies development, refactoring and testing.

Worker-pool basically is a bufferized channel of workers' input channels.
Once Worker is ready, it puts its input to the pool.
Once Manager is about to run a new task,
it takes a Worker input channel from the pool and send the Task config there.
Worker runs the Task, and when it's done, returns the input channel to the pool.

There are different syncronization techniques as well. Built both on channels and sync.WaitGroups.

For Tasks processing a TaskQueue is implemented.
It's a simple cycled queue. Just to implement some classic data structure.

The project has unit tests and e2e test.
The latter spins up a minio container to emulate S3, uploads a file there using S3UploadTask,
downloads the uploaded file with aws cli and compares checksums of both.

## demo
This will build a container, build an application inside it and run it with a set of simple 'sleep' tasks
    
    $ make demo