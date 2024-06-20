; ModuleID = 'build/example.ll'
source_filename = "example"

; Function Attrs: nofree nounwind
declare noundef i32 @printf(ptr nocapture noundef readonly, ...) local_unnamed_addr #0

; Function Attrs: nofree nosync nounwind memory(none)
define double @fib(double %n) local_unnamed_addr #1 {
entry:
  %equal = fcmp ueq double %n, 0.000000e+00
  br i1 %equal, label %common.ret, label %mergeBlock

common.ret:                                       ; preds = %mergeBlock, %entry, %mergeBlock4
  %common.ret.op = phi double [ %add, %mergeBlock4 ], [ 0.000000e+00, %entry ], [ 1.000000e+00, %mergeBlock ]
  ret double %common.ret.op

mergeBlock:                                       ; preds = %entry
  %equal1 = fcmp ueq double %n, 1.000000e+00
  br i1 %equal1, label %common.ret, label %mergeBlock4

mergeBlock4:                                      ; preds = %mergeBlock
  %subtract = fadd double %n, -1.000000e+00
  %callRes = tail call double @fib(double %subtract)
  %subtract5 = fadd double %n, -2.000000e+00
  %callRes6 = tail call double @fib(double %subtract5)
  %add = fadd double %callRes, %callRes6
  br label %common.ret
}

; Function Attrs: nofree nounwind
define void @main() local_unnamed_addr #0 {
entry:
  %callRes = tail call double @fib(double 4.800000e+01)
  %0 = alloca [18 x i8], align 1
  store i8 102, ptr %0, align 1
  %.fca.1.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 1
  store i8 105, ptr %.fca.1.gep, align 1
  %.fca.2.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 2
  store i8 98, ptr %.fca.2.gep, align 1
  %.fca.3.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 3
  store i8 40, ptr %.fca.3.gep, align 1
  %.fca.4.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 4
  store i8 37, ptr %.fca.4.gep, align 1
  %.fca.5.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 5
  store i8 46, ptr %.fca.5.gep, align 1
  %.fca.6.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 6
  store i8 48, ptr %.fca.6.gep, align 1
  %.fca.7.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 7
  store i8 102, ptr %.fca.7.gep, align 1
  %.fca.8.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 8
  store i8 41, ptr %.fca.8.gep, align 1
  %.fca.9.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 9
  store i8 32, ptr %.fca.9.gep, align 1
  %.fca.10.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 10
  store i8 105, ptr %.fca.10.gep, align 1
  %.fca.11.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 11
  store i8 115, ptr %.fca.11.gep, align 1
  %.fca.12.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 12
  store i8 32, ptr %.fca.12.gep, align 1
  %.fca.13.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 13
  store i8 37, ptr %.fca.13.gep, align 1
  %.fca.14.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 14
  store i8 46, ptr %.fca.14.gep, align 1
  %.fca.15.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 15
  store i8 48, ptr %.fca.15.gep, align 1
  %.fca.16.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 16
  store i8 102, ptr %.fca.16.gep, align 1
  %.fca.17.gep = getelementptr inbounds [18 x i8], ptr %0, i64 0, i64 17
  store i8 0, ptr %.fca.17.gep, align 1
  %callRes1 = call i32 (ptr, ...) @printf(ptr nonnull dereferenceable(1) %0, double 4.800000e+01, double %callRes)
  ret void
}

attributes #0 = { nofree nounwind }
attributes #1 = { nofree nosync nounwind memory(none) }
