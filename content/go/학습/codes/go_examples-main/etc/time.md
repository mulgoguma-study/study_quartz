Go의 time 패키지를 체계적으로 정리해드리겠습니다. 실무에서 자주 사용하는 패턴들을 중심으로 정리해볼게요.Go의 time 패키지를 실무에서 자주 사용하는 패턴들로 정리해봤습니다! 

**핵심 포인트:**

1. **포맷 레퍼런스 시간**: `2006-01-02 15:04:05` (1월 2일 3시 4분 5초 2006년)
2. **Parse vs ParseInLocation**: Parse는 UTC, ParseInLocation은 타임존 지정
3. **Add vs AddDate**: Add는 Duration, AddDate는 년/월/일 단위
4. **비교**: Before(), After(), Equal(), Compare() 사용
5. **타임존**: LoadLocation()으로 로드, In()으로 변환

**자주 하는 실수:**
- `time.Sleep(1000)` ❌ → `time.Sleep(1000 * time.Millisecond)` ✅
- 월은 1~12가 아니라 time.January ~ time.December 사용 권장
- Duration은 나노초 기반이므로 직접 정수 연산 주의

코드를 실행해보면서 각 기능들을 테스트해보세요. 특히 궁금한 부분이나 추가로 알고 싶은 패턴이 있으면 말씀해주세요!
