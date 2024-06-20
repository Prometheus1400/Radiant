; ModuleID = 'example'
source_filename = "example"

declare i32 @printf(ptr, ...)

define void @main() {
entry:
  %0 = alloca [13 x i8], align 1
  store [13 x i8] c"Hello Pookie\00", ptr %0, align 1
  %callRes = call i32 (ptr, ...) @printf(ptr %0)
  ret void
}
