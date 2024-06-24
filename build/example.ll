; ModuleID = 'example'
source_filename = "example"

declare i32 @printf(ptr, ...)

define void @main() {
entry:
  %x = alloca double, align 8
  store double 9.000000e+00, ptr %x, align 8
  %0 = load double, ptr %x, align 8
  %equal = fcmp ueq double %0, 1.000000e+00
  br i1 %equal, label %ifBlock, label %elseBlock

ifBlock:                                          ; preds = %entry
  %1 = alloca [7 x i8], align 1
  store [7 x i8] c"x is 1\00", ptr %1, align 1
  %callRes = call i32 (ptr, ...) @printf(ptr %1)
  br label %mergeBlock

elseBlock:                                        ; preds = %entry
  %2 = load double, ptr %x, align 8
  %equal1 = fcmp ueq double %2, 2.000000e+00
  br i1 %equal1, label %"elifBlock-%d0", label %"elifElseBlock-%d0"

mergeBlock:                                       ; preds = %"elifElseBlock-%d0", %"elifBlock-%d0", %ifBlock
  ret void

"elifBlock-%d0":                                  ; preds = %elseBlock
  %3 = alloca [7 x i8], align 1
  store [7 x i8] c"x is 2\00", ptr %3, align 1
  %callRes2 = call i32 (ptr, ...) @printf(ptr %3)
  br label %mergeBlock

"elifElseBlock-%d0":                              ; preds = %elseBlock
  %4 = alloca [7 x i8], align 1
  store [7 x i8] c"x is 0\00", ptr %4, align 1
  %callRes3 = call i32 (ptr, ...) @printf(ptr %4)
  br label %mergeBlock
}
