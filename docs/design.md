# Longhand — Design Outline

_A register VM in Go with bytecode you can read and write by hand._

## 1. Core commitments

These are the decisions that generate the rest of the system. Everything else should
fall out of them rather than be added alongside them.

1. **Registers are a window into one shared value stack.** A call frame is a base
   index, not an allocation. `R[i]` means `stack[base+i]`. Calls slide the base
   forward; returns slide it back. This is the single idea the whole design rests on.
2. **Fixed-width 32-bit instruction words**, with a small number of operand layouts
   (three-operand, one wide unsigned operand, one wide signed operand for jumps).
   Manual bit packing — the width limits are part of the lesson.
3. **Single goroutine, single interpreter loop.** No VM-level threads, no shared
   mutable heap across goroutines. Concurrency is explicitly out of scope.
4. **Go's GC initially**, via a tagged value type the Go runtime can trace. Writing
   your own collector is a later phase, taken on only once the instruction set has
   stopped moving.
5. **The text assembly is canonical**, not a debug view. Bytecode is a serialization
   of it, not the other way around.

---

## 2. Deliberate non-goals

Naming what you're _not_ building is what keeps the core clean.

| Not building                       | Why                                                                                                                                                                                     | Revisit when                                                                               |
| ---------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| VM-level concurrency               | Shared mutable heap across goroutines imports Go's data races into your language's semantics. There is no careful-interpreter-loop fix.                                                 | Never, in this project. It's a different project (actor VM, isolated heaps, copy-on-send). |
| Custom heap / own GC               | Requires giving up Go's tracing, which means giving up fast progress toward closures.                                                                                                   | Phase two, after upvalues work.                                                            |
| Multiple returns / varargs         | Forces a stack `top` that floats independently of the frame, and destroys the invariant that frame size is statically known. It is the largest single source of complexity in real Lua. | After fixed-arity calls are solid — and add it knowing what it costs.                      |
| RK operand encoding                | Stealing a high bit to mean "constant, not register" leaks into every decode site and burns register space. Lua 5.4 abandoned it.                                                       | Never. Use separate constant-operand opcodes instead.                                      |
| A source-language compiler         | It's a second project. Merging it in is how the clean core gets muddy.                                                                                                                  | After 1.0 is frozen — as a separate repo. See section 7.                                   |
| JIT, NaN-boxing, superinstructions | Optimizations that presuppose a stable design.                                                                                                                                          | Not in scope.                                                                              |

---

## 3. Build order

Each phase should be independently finishable and testable. Resist starting the next
one early.

**Phase 0 — Instruction encoding + assembler/disassembler.**
Before any execution. The key discipline: property-test the round trip so that
assembling and disassembling reaches a fixpoint. This catches encoding bugs long
before the interpreter would, and it forces the text format to be _complete_ rather
than lossy.

**Phase 1 — Value representation and the stack.**
A tagged value the Go GC can trace. Get the register-window arithmetic right here;
it's cheap to fix now and expensive later.

**Phase 2 — Interpreter loop, arithmetic, jumps.**
Straight-line and branching code. No calls yet.

**Phase 3 — Fixed-arity calls and frames.**
Base sliding, return values, call depth limits. This is where phase 1's design
either pays off or doesn't.

**Phase 4 — Heap objects.**
Strings, and a table/map type if you want one. Interning decisions belong here.

**Phase 5 — Closures and upvalues.**
The real prize. Open upvalues point directly into the live stack; closing them on
scope exit copies the value into a heap cell and repoints the reference. This is the
mechanism that lets closures work without boxing every local, and implementing it is
the moment the whole design clicks.

**Phase 6 (optional) — Replace Go's GC with your own.**
Only once the instruction set is frozen. Mark-sweep first; generational or
incremental only if you still care afterward.

---

## 4. Questions to answer explicitly before coding each phase

Write down the answer. Don't discover it mid-implementation.

- How many registers can one operand field address, and what happens when a function
  needs more?
- Where do constants live — per-function pool, or global? How wide is the index?
- What is the exact protocol for a call: who writes the arguments, who moves the
  return value, who restores the base?
- Is your dispatch loop a `switch`, a function table, or computed goto emulation? What
  does each cost you in Go specifically?
- How does the assembler resolve forward jump labels, and what's the max jump distance?
- What's the error model — trap and unwind, or a status value threaded through?
- What exactly does the disassembler need to emit for the round trip to be lossless?

---

## 5. Where the time should actually go

Effort should be wildly unevenly distributed. The useful sorting question is not
"what's hard?" but **"what is expensive to change later?"** Those are different sets.

**Spend the bulk of your time here (~60%):**

- **Value representation and the register window (Phase 1).** Everything downstream
  reads and writes through this. Getting the base-index arithmetic and the value type
  wrong means rewriting the interpreter, not patching it. Think about this before you
  type, and be willing to throw away a first attempt.
- **The calling protocol (Phase 3).** Who writes arguments, who moves results, who
  restores the base — this is a contract every future opcode depends on. Write it down
  in prose before implementing. Ambiguity here shows up as heisenbugs much later.
- **Upvalues (Phase 5).** The one place where the code is short and the thinking is
  long. Expect to spend days understanding open-vs-closed and the close-on-scope-exit
  ordering. This is the highest learning-per-line ratio in the entire project.

**Spend real but bounded time here (~25%):**

- **Tooling: a tracer and a state dumper.** Underrated, and the thing most people skip.
  Debugging a VM without the ability to single-step and print the stack, frame bases,
  and open upvalue list is genuinely miserable — you will lose more hours to its
  absence than you would ever have spent building it. Build it during Phase 2, not when
  you're desperate in Phase 5.
- **The round-trip property test.** Cheap to write, catches an entire class of bug for
  the rest of the project's life.
- **Encoding and the assembler (Phase 0).** Fiddly but shallow. It should feel like
  careful clerical work, not like a research problem. If it's taking more than a few
  days, you're over-designing the format.

**Deliberately underspend here (~15%):**

- **The dispatch loop.** It's a `switch`. Write it in an afternoon. Do not benchmark it,
  do not attempt threaded dispatch, do not read about computed goto until the VM works.
  It is also the cheapest part to rewrite later, which is precisely why it doesn't
  deserve up-front agonizing.
- **Opcode breadth.** Adding the 30th arithmetic variant _feels_ like progress and
  teaches nothing. A small instruction set that supports closures beats a large one
  that doesn't.
- **Assembly syntax bikeshedding.** Pick a syntax in an hour. The round-trip test is
  what matters; the surface is not.

The failure mode to watch for: the underspend list is more immediately gratifying than
the overspend list. Dispatch tuning and new opcodes produce a steady drip of visible
wins. Register-window design produces a day of staring at a whiteboard. Notice when
you're reaching for the former to avoid the latter.

---

## 6. Motivating mechanisms that lead to demoable moments

1. First arithmetic program runs from hand-written assembly.
2. First recursive call — naive fibonacci — proving frames and the base slide work.
3. First closure that captures a mutable counter and outlives its scope. This is the
   one that will feel like magic, and it's why the phase order puts it last.
4. First execution trace that a stranger can read and follow without your help.

---

## 7. After 1.0 — front-ends as validation

A **separate project**, started only once 1.0 is frozen. Its purpose is to generate the
change list for v2, not to ship languages. Two front-ends, two distinct jobs.

### 7.1 Track A — Scheme (design pressure)

The job: make the VM answer questions it can't answer yet.

- Cheapest possible front-end. The reader _is_ the parser; there is no grammar to write.
- **Proper tail calls are the point.** A tail call must reuse the current frame's base
  rather than pushing a new one. This is a hard, testable demand on the calling
  protocol from section 3, and it is the single most likely thing to change the VM.
- **Optional type annotations, deliberately dumb.** Trusted and unchecked — a codegen
  hint only. If the program lies, behaviour is undefined. This buys the entire lesson
  worth having (typed opcode variants, unboxed representations, elided tag checks) for
  a fraction of the cost of a real type system.
- **Guard against the drift into a type system.** The moment you add a checker, you
  have soundness obligations, inference, error messages, and a second project. If you
  ever want checking, it is strictly additive later. Write this constraint down.
- Optional stretch, only if the appetite is there: continuations. Genuinely
  informative about frame ownership, and genuinely expensive.

Dropped from an earlier version of this plan, and why: **Haskell** (laziness wants
graph reduction, not a register machine — the lesson is a weekend reading about STG,
not a year implementing it), **Python** (generators contradict "a frame is a base index"
on day one — a v2 decision to make deliberately, not to stumble into), **JavaScript**
(VM-level lessons overlap almost entirely with Scheme's, while the grammar is hostile),
and **full ML** (pattern matching is a compiler problem, not a VM one; the only parts
that reached the VM were ADT representation and typed opcodes, and annotated Scheme
delivers both).

### 7.2 Track B — Lua subset (conformance and regression)

The job: a differential-testing oracle against the reference Lua VM.

- **Know what this test can and cannot do.** Compiling Lua onto a Lua-derived VM is a
  _conformance_ test, not a falsification test. It will find real bugs. It cannot tell
  you the design is wrong, because the design was taken from the thing producing the
  workload. Don't let it produce false confidence.
- **Pin the version** (5.4 or 5.3 — they differ on integers) and **write the subset
  down as a spec before writing the compiler.** Without that, every mismatch becomes
  "bug, or just outside my subset?" and the oracle degrades into noise.
- Explicit exclusions to state up front: metatables and `__index` chains, string↔number
  coercion, the integer/float distinction, `goto`, weak tables.
- **Compare observable output only — never internals.** The moment you are diffing
  stack layouts against reference Lua, you have coupled your VM to theirs and lost the
  freedom to change it.
- This is where multiple returns and coroutines get exercised honestly, if 1.0 has
  grown them.

### 7.3 Sequencing

**Scheme first, then Lua.** Tail calls will likely change the calling protocol; building
a Lua conformance corpus first means invalidating it. Settle the VM's shape under design
pressure, _then_ lay down the regression net.

**Serially, with VM fixes in between** — not both against a frozen 1.0. Two front-ends
built in parallel produce two piles of complaints with no way to attribute which fix
mattered. Compile → find the gap → fix the VM → re-verify is what makes this a feedback
loop rather than a survey.

The quiet payoff, and the real one: once Track B exists, you have a fast, unambiguous
signal for whether a v2 change broke something. That makes you willing to make
_aggressive_ v2 changes. That matters more than any individual bug either track finds.

---

## 8. Reading

### Primary sources

- _The Implementation of Lua 5.0_ — Ierusalimschy, de Figueiredo, Celes. The single
  most useful paper for this project. Covers register allocation, the value stack,
  upvalues, and why they went register-based. `https://www.lua.org/doc/jucs05.pdf`
- Lua papers index: `https://www.lua.org/papers.html`
- Lua 5.4 source, especially `lopcodes.h`, `lvm.c`, `lfunc.c`, `ldo.c`.
  `https://github.com/lua/lua`
- _Crafting Interpreters_, Robert Nystrom — Part III builds a stack VM in C, but the
  chapters on closures, garbage collection, and value representation transfer directly.
  `https://craftinginterpreters.com`

### Comparison implementations worth reading

- **Wren** (`https://github.com/wren-lang/wren`) — small, clean, well-commented, by the
  author of _Crafting Interpreters_.
- **LuaJIT wiki / mailing list** — Mike Pall's writing on interpreter dispatch and
  bytecode design is exceptional even if you never touch JIT.
- **CPython's `ceval.c`** — a counterexample. Stack-based, and instructive precisely
  for the contrast.
- **Dalvik / Android bytecode spec** — another register machine, with very different
  encoding tradeoffs from Lua's.

### Search terms

Design and encoding:
`register-based virtual machine design` ·
`bytecode instruction encoding operand width` ·
`three-address code bytecode` ·
`register window calling convention` ·
`Lua register allocation compiler`

Interpreter mechanics:
`interpreter dispatch techniques` ·
`direct threaded code vs switch dispatch` ·
`computed goto interpreter` ·
`Go interpreter loop bounds check elimination` ·
`instruction decoding branch prediction interpreter`

Closures and memory:
`upvalue open closed closure implementation` ·
`flat closure conversion` ·
`NaN boxing value representation` (background — you're deferring this) ·
`tagged union dynamic value representation` ·
`mark sweep garbage collector implementation` ·
`precise vs conservative garbage collection roots`

Go-specific:
`Go interface dispatch cost` ·
`Go escape analysis heap allocation` ·
`unsafe.Pointer garbage collector rules` (needed only if you reach phase 6) ·
`Go slice bounds check elimination`

Comparative:
`stack vs register virtual machine performance` — look for Shi et al., _Virtual Machine
Showdown: Stack Versus Registers_, which is the paper most of this debate cites.

Front-end tracks (section 7, not needed until after 1.0):
`proper tail call implementation` ·
`tail call optimization frame reuse interpreter` ·
`R7RS small standard` ·
`Scheme from Scratch` / `write yourself a Scheme` ·
`differential testing compiler` ·
`type-specialized bytecode opcodes` ·
`unboxed value representation dynamic language` ·
`algebraic data type memory representation tagged` ·
`STG machine lazy evaluation` (background only — for understanding why laziness was
excluded, not for implementing it)

---

## 9. The main risk

Not difficulty — scope. The three tempting expansions are concurrency, multiple
returns, and a source language, and each one is individually reasonable and
collectively fatal to a clean core. If you find yourself adding one before the phase
list above is finished, that's the signal to stop rather than the signal to push on.
