## Pointers

Pointers provide a way to share data across function boundaries. Having the ability to share and reference data with a pointer provides flexbility. It also helps our programs minimize the amount of memory they need and can add some extra performance.

## Notes

* Use pointers to share data.
* Values in Go are always pass by value.
* "Value of", what's in the box. "Address of" ( **&** ), where is the box.
* The (*) operator declares a pointer variable and the "Value that the pointer points to".

## Garbage Collection

The design of the Go GC has changed over the years:
* Go 1.0, Stop the world mark sweep collector based heavily on tcmalloc.
* Go 1.2, Precise collector, wouldn't mistake big numbers (or big strings of text) for pointers.
* Go 1.3, Fully precise tracking of all stack values.
* Go 1.4, Mark and sweep now parallel, but still stop the world.
* Go 1.5, New GC design, focusing on latency over throughput.
* Go 1.6, GC improvements, handling larger heaps with lower latency.
* Go 1.7, GC improvements, handling larger number of idle goroutines, substantial stack size fluctuation, or large package-level variables.
* Go 1.8, GC improvements, collection pauses should be significantly shorter than they were in Go 1.7, usually under 100 microseconds and often as low as 10 microseconds.

![figure1](GC_Algorithm.png?v=4)

**STW : Stop The World Phase**

Turn the Write Barrier on. The Write Barrier is a little function that inspects the write of pointers when the GC is running. Each goroutine must know this flag is set. This STW pause should be sub-millisecond.

**Mark Phase**

Find all the objects that can be reclaimed.

* All objects on the heap are turned WHITE.
* **Scan Stacks :** Find all the root objects and place them in the queue.
    * Pause the goroutine while scanning its stack.
    * All root objects found on the stack are turned GREY.
    * The stack is marked BLACK.
* **Mark Phase I :** Pop a GREY object from the queue and scan it.
    * Turn the object BLACK.
    * If this BLACK object points to a WHITE object, the WHITE object is turned GREY and added to the queue.
    * The GC and the application are running concurrently.
    * Goroutines executing at this time will find their stack reverted back to GREY.
* **Mark Phase II - STW :** Re-scan GREY stacks.
    * Re-scan all GREY stacks and root objects again.
    * Should be quick but large numbers of active goroutines can cause milliseconds of latency. 
    * Call any finalizers on BLACK objects.

**Sweep Phase**

Sweep phase reclaims memory.

* Left with either WHITE or BLACK objects. No more GREY objects.
* WHITE objects are reclaimed while BLACK objects stay.

**Write Barrier**

The Write Barrier is a little function that inspects the write of pointers when the GC is running.

The Write Barrier is in place to prevent a situation where a BLACK object (one that is processed) suddenly finds itself pointing to a WHITE object after the Mark Phases are complete. This could happen if a goroutine changes (writes) a pointer inside a BLACK object to point to a WHITE object while both the program and the GC is running after that BLACK object has been processed. So the Write Barrier will make sure this write changes the object to BLACK so it's not swept away.

Pointers to the heap that exist on a stack can also be changed by goroutines when the GC is running. So stacks are also marked as BLACK once they are scanned and can revert back to GREY during Mark Phase I. A BLACK stack reverts back to GREY when its goroutine executes again. During Mark Phase II, the GC must re-scan GREY stacks to BLACKen them and finish marking any remaining heap pointers. Since it must ensure the stacks don't continue to change during this scan, the whole re-scan process happens *with the world stopped*.

**Pacing**

The GC starts a scan based on a feedback loop of information about the running program and the stress on the heap. It is the pacers job to make this decision. Once the decision is made to run, the amount of time the GC has to finish the scan is pre-determined. This time is based on the current size of the heap, the current live heap, and timing calculations about future heap usage while the GC is running.

The GC has a set of goroutines to perform the task of Mark and Sweep. The scheduler will provide these goroutines 25% of the available logical processor time. If your program is using 4 logical processors, that 1 entire logical processor will be given to the GC goroutines for exclusive use.

If the GC begins to believe that it can’t finish the collection within the decided amount of time, it will begin to recruit program goroutines to help. Those goroutines that are causing the slow down will be recruited to help.

## Links

### Pointer Mechanics

https://golang.org/doc/effective_go.html#pointers_vs_values  
https://www.goinggo.net/2017/05/language-mechanics-on-stacks-and-pointers.html  
http://www.goinggo.net/2014/12/using-pointers-in-go.html  
http://www.goinggo.net/2013/07/understanding-pointers-and-memory.html  

### Stacks

[Contiguous Stack Proposal](https://docs.google.com/document/d/1wAaf1rYoM4S4gtnPh0zOlGzWtrZFQ5suE8qr2sD8uWQ/pub)  

### Escape Analysis and Inlining

[Go Escape Analysis Flaws](https://docs.google.com/document/d/1CxgUBPlx9iJzkz9JWkb6tIpTe5q32QDmz8l0BouG0Cw)  
[Compiler Optimizations](https://github.com/golang/go/wiki/CompilerOptimizations)

### Garbage Collection

[The Garbage Collection Handbook](http://gchandbook.org/)  
[Tracing Garbage Collection](https://en.wikipedia.org/wiki/Tracing_garbage_collection)  
[Go Blog - 1.5 GC](https://blog.golang.org/go15gc)  
[Go GC: Solving the Latency Problem](https://www.youtube.com/watch?v=aiv1JOfMjm0&index=16&list=PL2ntRZ1ySWBf-_z-gHCOR2N156Nw930Hm)  
[Concurrent garbage collection](http://rubinius.com/2013/06/22/concurrent-garbage-collection)  
[Go 1.5 concurrent garbage collector pacing](https://docs.google.com/document/d/1wmjrocXIWTr1JxU-3EQBI6BK6KgtiFArkG47XK73xIQ/edit)  
[Eliminating Stack Re-Scanning](https://github.com/golang/proposal/blob/master/design/17503-eliminate-rescan.md)  
[Why golang garbage-collector not implement Generational and Compact gc?](https://groups.google.com/forum/m/#!topic/golang-nuts/KJiyv2mV2pU) - Ian Lance Taylor  

### Static Single Assignment Optimizations

[GopherCon 2015: Ben Johnson - Static Code Analysis Using SSA](https://www.youtube.com/watch?v=D2-gaMvWfQY)  
https://godoc.org/golang.org/x/tools/go/ssa  
[Understanding Compiler Optimization](https://www.youtube.com/watch?v=FnGCDLhaxKU)

### Debugging code generation

[Debugging code generation in Go](http://golang.rakyll.org/codegen/)

## Code Review

[Pass by Value](example1/example1.go) ([Go Playground](https://play.golang.org/p/9kxh18hd_BT))  
[Sharing data I](example2/example2.go) ([Go Playground](https://play.golang.org/p/mJz5RINaimn))  
[Sharing data II](example3/example3.go) ([Go Playground](https://play.golang.org/p/GpmPICMGMre))  
[Escape Analysis](example4/example4.go) ([Go Playground](https://play.golang.org/p/qNOw5gEtYhI))  
[Stack grow](example5/example5.go) ([Go Playground](https://play.golang.org/p/BgIKcFcZ4PO))  

### Escape Analysis Flaws

[Indirect Assignment](flaws/example1/example1_test.go)  
[Indirection Execution](flaws/example2/example2_test.go)  
[Assignment Slices Maps](flaws/example3/example3_test.go)  
[Indirection Level Interfaces](flaws/example4/example4_test.go)  
[Unknown](flaws/example5/example5_test.go)  

## Exercises

### Exercise 1

**Part A** Declare and initialize a variable of type int with the value of 20. Display the _address of_ and _value of_ the variable.

**Part B** Declare and initialize a pointer variable of type int that points to the last variable you just created. Display the _address of_ , _value of_ and the _value that the pointer points to_.

[Template](exercises/template1/template1.go) ([Go Playground](https://play.golang.org/p/6QYTKWzF8s8)) |
[Answer](exercises/exercise1/exercise1.go) ([Go Playground](https://play.golang.org/p/qq5P9gRDHKc))

### Exercise 2

Declare a struct type and create a value of this type. Declare a function that can change the value of some field in this struct type. Display the value before and after the call to your function.

[Template](exercises/template2/template2.go) ([Go Playground](https://play.golang.org/p/nolKjrgBX-X)) |
[Answer](exercises/exercise2/exercise2.go) ([Go Playground](https://play.golang.org/p/i6utWhgDUH4))
___
All material is licensed under the [Apache License Version 2.0, January 2004](http://www.apache.org/licenses/LICENSE-2.0).
