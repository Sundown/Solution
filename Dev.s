	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 12, 0
	.globl	_double                         ## -- Begin function double
	.p2align	4, 0x90
_double:                                ## @double
	.cfi_startproc
## %bb.0:                               ## %entry
	leaq	(%rdi,%rdi), %rax
	retq
	.cfi_endproc
                                        ## -- End function
	.globl	_Start                          ## -- Begin function Start
	.p2align	4, 0x90
_Start:                                 ## @Start
	.cfi_startproc
## %bb.0:                               ## %entry
	subq	$24, %rsp
	.cfi_def_cfa_offset 32
	movl	$4, %edi
	movl	$8, %esi
	callq	_calloc
	movq	$1, (%rsp)
	movq	$8, 8(%rsp)
	movq	%rax, 16(%rsp)
	movq	$1, (%rax)
	movq	(%rsp), %rax
	movq	8(%rsp), %rdx
	movq	16(%rsp), %rcx
	addq	$24, %rsp
	retq
	.cfi_endproc
                                        ## -- End function
	.globl	_main                           ## -- Begin function main
	.p2align	4, 0x90
_main:                                  ## @main
	.cfi_startproc
## %bb.0:                               ## %entry
	pushq	%rax
	.cfi_def_cfa_offset 16
	callq	_Start
	xorl	%eax, %eax
	popq	%rcx
	retq
	.cfi_endproc
                                        ## -- End function
.subsections_via_symbols
