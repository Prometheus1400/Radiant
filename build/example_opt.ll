; ModuleID = 'build/example.ll'
source_filename = "example"

declare i32 @printf(ptr, ...)

define double @fib(double %n) {
entry:
  %equal = fcmp ueq double %n, 0.000000e+00
  br i1 %equal, label %ifBlock, label %elseBlock

ifBlock:                                          ; preds = %entry
  ret double 0.000000e+00

0:                                                ; No predecessors!
  br label %mergeBlock

elseBlock:                                        ; preds = %entry
  br label %mergeBlock

mergeBlock:                                       ; preds = %elseBlock, %0
  %equal1 = fcmp ueq double %n, 1.000000e+00
  br i1 %equal1, label %ifBlock2, label %elseBlock3

ifBlock2:                                         ; preds = %mergeBlock
  ret double 1.000000e+00

1:                                                ; No predecessors!
  br label %mergeBlock4

elseBlock3:                                       ; preds = %mergeBlock
  br label %mergeBlock4

mergeBlock4:                                      ; preds = %elseBlock3, %1
  %subtract = fsub double %n, 1.000000e+00
  %callRes = call double @fib(double %subtract)
  %subtract5 = fsub double %n, 2.000000e+00
  %callRes6 = call double @fib(double %subtract5)
  %add = fadd double %callRes, %callRes6
  ret double %add
}

define void @main() {
entry:
  %num = alloca double, align 8
  store double 4.000000e+01, ptr %num, align 8
  %0 = load double, ptr %num, align 8
  %callRes = call double @fib(double %0)
  %res = alloca double, align 8
  store double %callRes, ptr %res, align 8
  %1 = alloca [18 x i8], align 1
  store [18 x i8] c"fib(%.0f) is %.0f\00", ptr %1, align 1
  %2 = load double, ptr %num, align 8
  %3 = load double, ptr %res, align 8
  %callRes1 = call i32 (ptr, ...) @printf(ptr %1, double %2, double %3)
  ret void
}
