Context Type, Carrier deadlines, cancel signal
그리고 기타 request-scoped values를  API와 프로세스, 프로세스와 프로세스간 전송하는 데이터

파생된 전파방법은 WithCancel, WithDeadline, WithTimeout, WithValue.

WithCancelCasue, WithDeadlineCasue, WithTimeoutCause, WithValueCause 와 같은 방식으로 error 원인도 리턴가능. return ctx.Err() 와 같은 방식

