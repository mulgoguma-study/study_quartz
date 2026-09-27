package huhu1

/*
================================================================================
문제 4: 자율주행 차량 상태 머신 - 솔루션
================================================================================

【 핵심 설계 원칙 】
1. 전이 테이블 기반 - O(1) 상태 전이
2. Guard/Action 패턴 - 유연한 조건과 액션
3. Observer 패턴 - 상태 진입/퇴장 콜백
4. Ring Buffer - 고정 크기 히스토리

【 상태 전이도 】
┌─────────┐  START   ┌─────────┐  ACCELERATE  ┌──────────┐
│ PARKED  │───────→  │  IDLE   │────────────→ │ CRUISING │
└─────────┘          └─────────┘              └──────────┘
     ↑                    ↑                        │
     │ STOP               │ BRAKE                  │
     │                    └────────────────────────┘
     │                                             │
     │                                    LANE_CHANGE_START
     │                                             ↓
     │                                    ┌──────────────┐
     │←───────────────── EMERGENCY ───────│LANE_CHANGING │
     │                    (from any)      └──────────────┘
     ↓
┌────────────────┐
│ EMERGENCY_STOP │
└────────────────┘

================================================================================
*/

import (
	"errors"
	"sync"
	"time"
)

// ============================================================================
// 전이 키 (상태, 이벤트 조합)
// ============================================================================

type transitionKey struct {
	from  VehicleState
	event VehicleEvent
}

// ============================================================================
// 링 버퍼 (히스토리용)
// ============================================================================

type historyRingBuffer struct {
	data  []StateHistory
	head  int
	size  int
	count int
}

func newHistoryRingBuffer(size int) *historyRingBuffer {
	return &historyRingBuffer{
		data: make([]StateHistory, size),
		size: size,
	}
}

func (r *historyRingBuffer) Push(h StateHistory) {
	r.data[r.head] = h
	r.head = (r.head + 1) % r.size
	if r.count < r.size {
		r.count++
	}
}

func (r *historyRingBuffer) GetAll() []StateHistory {
	result := make([]StateHistory, r.count)

	if r.count < r.size {
		// 아직 버퍼가 가득 차지 않음
		copy(result, r.data[:r.count])
	} else {
		// 버퍼가 가득 참 - head부터 순환하여 복사
		idx := 0
		for i := r.head; idx < r.size; i = (i + 1) % r.size {
			result[idx] = r.data[i]
			idx++
		}
	}

	return result
}

// ============================================================================
// FSM 구현
// ============================================================================

// VehicleFSMSolution 완성된 상태 머신 구현
type VehicleFSMSolution struct {
	mu sync.RWMutex

	// 현재 상태
	currentState VehicleState

	// 전이 테이블: (from, event) -> transition
	transitions map[transitionKey]*StateTransition

	// 콜백
	onEnterHandlers map[VehicleState][]func()
	onExitHandlers  map[VehicleState][]func()

	// 히스토리
	history *historyRingBuffer
}

// NewVehicleFSMSolution 생성자
func NewVehicleFSMSolution(initial VehicleState) *VehicleFSMSolution {
	fsm := &VehicleFSMSolution{
		currentState:    initial,
		transitions:     make(map[transitionKey]*StateTransition),
		onEnterHandlers: make(map[VehicleState][]func()),
		onExitHandlers:  make(map[VehicleState][]func()),
		history:         newHistoryRingBuffer(100),
	}

	// 기본 전이 규칙 등록
	fsm.registerDefaultTransitions()

	return fsm
}

// registerDefaultTransitions 기본 전이 규칙 등록
func (f *VehicleFSMSolution) registerDefaultTransitions() {
	// PARKED 상태에서의 전이
	f.RegisterTransition(StateTransition{
		From:  PARKED,
		Event: START,
		To:    IDLE,
	})
	f.RegisterTransition(StateTransition{
		From:  PARKED,
		Event: CHARGE_START,
		To:    CHARGING,
	})

	// IDLE 상태에서의 전이
	f.RegisterTransition(StateTransition{
		From:  IDLE,
		Event: STOP,
		To:    PARKED,
	})
	f.RegisterTransition(StateTransition{
		From:  IDLE,
		Event: ACCELERATE,
		To:    CRUISING,
	})

	// CRUISING 상태에서의 전이
	f.RegisterTransition(StateTransition{
		From:  CRUISING,
		Event: BRAKE,
		To:    IDLE,
	})
	f.RegisterTransition(StateTransition{
		From:  CRUISING,
		Event: LANE_CHANGE_START,
		To:    LANE_CHANGING,
	})
	f.RegisterTransition(StateTransition{
		From:  CRUISING,
		Event: TURN_START,
		To:    TURNING,
	})

	// LANE_CHANGING 상태에서의 전이
	f.RegisterTransition(StateTransition{
		From:  LANE_CHANGING,
		Event: LANE_CHANGE_COMPLETE,
		To:    CRUISING,
	})

	// TURNING 상태에서의 전이
	f.RegisterTransition(StateTransition{
		From:  TURNING,
		Event: TURN_COMPLETE,
		To:    CRUISING,
	})

	// EMERGENCY_STOP 상태에서의 전이
	f.RegisterTransition(StateTransition{
		From:  EMERGENCY_STOP,
		Event: RESUME,
		To:    IDLE,
	})

	// CHARGING 상태에서의 전이
	f.RegisterTransition(StateTransition{
		From:  CHARGING,
		Event: CHARGE_COMPLETE,
		To:    PARKED,
	})

	// 모든 상태에서 EMERGENCY 이벤트 → EMERGENCY_STOP
	for _, state := range []VehicleState{IDLE, CRUISING, LANE_CHANGING, TURNING} {
		f.RegisterTransition(StateTransition{
			From:  state,
			Event: EVENT_EMERGENCY,
			To:    EMERGENCY_STOP,
		})
	}
}

// RegisterTransition 전이 규칙 등록
func (f *VehicleFSMSolution) RegisterTransition(t StateTransition) {
	f.mu.Lock()
	defer f.mu.Unlock()

	key := transitionKey{from: t.From, event: t.Event}
	f.transitions[key] = &t
}

// GetCurrentState 현재 상태 반환
func (f *VehicleFSMSolution) GetCurrentState() VehicleState {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.currentState
}

// Transition 상태 전이 실행
/*
【 시간 복잡도 】O(1) + O(콜백 개수)
【 동작 원리 】
1. 전이 테이블에서 전이 규칙 조회
2. Guard 조건 확인 (있으면)
3. 현재 상태의 OnExit 콜백 실행
4. 전이 Action 실행 (있으면)
5. 상태 변경
6. 새 상태의 OnEnter 콜백 실행
7. 히스토리에 기록
*/
func (f *VehicleFSMSolution) Transition(event VehicleEvent) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 1. 전이 규칙 조회
	key := transitionKey{from: f.currentState, event: event}
	transition, ok := f.transitions[key]
	if !ok {
		return errors.New("invalid transition: no rule for (" + string(f.currentState) + ", " + string(event) + ")")
	}

	// 2. Guard 조건 확인
	if transition.Guard != nil && !transition.Guard() {
		return errors.New("transition blocked by guard condition")
	}

	oldState := f.currentState
	newState := transition.To

	// 3. OnExit 콜백
	if handlers, ok := f.onExitHandlers[oldState]; ok {
		for _, handler := range handlers {
			handler()
		}
	}

	// 4. Action 실행
	if transition.Action != nil {
		transition.Action()
	}

	// 5. 상태 변경
	f.currentState = newState

	// 6. OnEnter 콜백
	if handlers, ok := f.onEnterHandlers[newState]; ok {
		for _, handler := range handlers {
			handler()
		}
	}

	// 7. 히스토리 기록
	f.history.Push(StateHistory{
		From:      oldState,
		To:        newState,
		Event:     event,
		Timestamp: time.Now(),
	})

	return nil
}

// CanTransition 전이 가능 여부 확인
func (f *VehicleFSMSolution) CanTransition(event VehicleEvent) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	key := transitionKey{from: f.currentState, event: event}
	transition, ok := f.transitions[key]
	if !ok {
		return false
	}

	// Guard 조건 확인
	if transition.Guard != nil {
		return transition.Guard()
	}

	return true
}

// OnEnter 상태 진입 콜백 등록
func (f *VehicleFSMSolution) OnEnter(state VehicleState, handler func()) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.onEnterHandlers[state] = append(f.onEnterHandlers[state], handler)
}

// OnExit 상태 퇴장 콜백 등록
func (f *VehicleFSMSolution) OnExit(state VehicleState, handler func()) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.onExitHandlers[state] = append(f.onExitHandlers[state], handler)
}

// GetStateHistory 상태 변경 이력 반환
func (f *VehicleFSMSolution) GetStateHistory() []StateHistory {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.history.GetAll()
}

// ============================================================================
// 고급 기능: 계층적 상태 머신 (HSM)
// ============================================================================

/*
【 Hierarchical State Machine 】
복잡한 자율주행 시스템에서는 계층적 상태 머신이 필요합니다.

예시:
DRIVING (상위 상태)
├── CRUISING (하위)
├── LANE_CHANGING (하위)
└── TURNING (하위)

장점:
- 공통 동작을 상위 상태에서 정의
- 상태 폭발 문제 해결
- 코드 재사용성 향상

구현 힌트:
type HierarchicalState struct {
    Name     VehicleState
    Parent   *HierarchicalState
    Children []*HierarchicalState
}

// 현재 상태 또는 부모에서 전이 규칙 찾기
func (f *FSM) findTransition(state *HierarchicalState, event VehicleEvent) *Transition {
    for s := state; s != nil; s = s.Parent {
        key := transitionKey{from: s.Name, event: event}
        if t, ok := f.transitions[key]; ok {
            return t
        }
    }
    return nil
}
*/

// ============================================================================
// 테스트용 메인 함수
// ============================================================================

func mainFSM() {
	// FSM 생성
	fsm := NewVehicleFSMSolution(PARKED)

	// 콜백 등록
	fsm.OnEnter(CRUISING, func() {
		println("Entered CRUISING state - activating cruise control")
	})

	fsm.OnEnter(EMERGENCY_STOP, func() {
		println("EMERGENCY STOP - activating hazard lights")
	})

	fsm.OnExit(CRUISING, func() {
		println("Exiting CRUISING state - deactivating cruise control")
	})

	// 상태 전이 시뮬레이션
	println("Initial state:", string(fsm.GetCurrentState()))

	// PARKED → IDLE
	if err := fsm.Transition(START); err != nil {
		println("Error:", err.Error())
	}
	println("After START:", string(fsm.GetCurrentState()))

	// IDLE → CRUISING
	if err := fsm.Transition(ACCELERATE); err != nil {
		println("Error:", err.Error())
	}
	println("After ACCELERATE:", string(fsm.GetCurrentState()))

	// CRUISING → LANE_CHANGING
	if err := fsm.Transition(LANE_CHANGE_START); err != nil {
		println("Error:", err.Error())
	}
	println("After LANE_CHANGE_START:", string(fsm.GetCurrentState()))

	// 비상 정지!
	if err := fsm.Transition(EVENT_EMERGENCY); err != nil {
		println("Error:", err.Error())
	}
	println("After EMERGENCY:", string(fsm.GetCurrentState()))

	// 히스토리 출력
	println("\nState History:")
	for _, h := range fsm.GetStateHistory() {
		println(string(h.From), "→", string(h.To), "via", string(h.Event))
	}
}

// ============================================================================
// 성능 최적화 포인트
// ============================================================================

/*
【 프로덕션 최적화 】

1. Lock-free 읽기
   - atomic.Value로 현재 상태 저장
   - 읽기 시 락 불필요

2. 이벤트 기반 아키텍처
   - 채널로 이벤트 수신
   - 전용 고루틴에서 순차 처리
   - 동시성 이슈 원천 차단

3. 상태 직렬화
   - 상태 머신 스냅샷 저장
   - 장애 시 복구 가능

4. 메트릭 수집
   - 상태별 체류 시간
   - 전이 빈도
   - 실패한 전이 수

┌─────────────────────────────────────────────────────────────────────────────┐
│                       프로덕션 아키텍처                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  [센서 데이터] → [Event Generator] → [Event Channel] → [FSM Goroutine]     │
│                                                              ↓              │
│                                                    [State Persistence]      │
│                                                              ↓              │
│  [Actuator Control] ← [Action Executor] ← [State Change Notification]     │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

【 실제 자율주행 상태 예시 】

Level 4 자율주행 상태:
- MANUAL: 수동 운전
- AUTOMATED_READY: 자율주행 대기
- AUTOMATED_DRIVING: 자율주행 중
  - LANE_FOLLOWING: 차선 유지
  - LANE_CHANGE: 차선 변경
  - INTERSECTION: 교차로 통과
  - PARKING: 주차 중
- TAKEOVER_REQUEST: 운전자 제어권 인수 요청
- MINIMAL_RISK_CONDITION: 최소 위험 상태 (안전하게 정차)
*/
