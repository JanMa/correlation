package correlation

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/uuid"
)

var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("bar"))
})

func TestNoConfig(t *testing.T) {
	c := New(Options{})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	c.Handler(handler).ServeHTTP(res, req)

	expectEq(t, res.Code, http.StatusOK)
	expectEq(t, res.Body.String(), "bar")
	expectNeq(t, res.Header().Get(correlationIDHeader), "")
}

func TestRequestOnly(t *testing.T) {
	c := New(Options{})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	c.HandlerForRequestOnly(handler).ServeHTTP(res, req)

	expectEq(t, res.Code, http.StatusOK)
	if _, ok := uuid.Parse(req.Header.Get(correlationIDHeader)); ok != nil {
		t.Errorf("Expected valid UUID - Got [%v]. Error: [%v]", req.Header.Get(correlationIDHeader), ok)
	}
	expectEq(t, res.Header().Get(correlationIDHeader), "")

}

func TestHeaderName(t *testing.T) {
	c := New(Options{
		HeaderName: "Foo",
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	c.Handler(handler).ServeHTTP(res, req)
	expectNeq(t, res.Header().Get(c.opt.HeaderName), "")
}

func TestUUIDHeader(t *testing.T) {
	c := New(Options{
		IDType: UUID,
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	c.Handler(handler).ServeHTTP(res, req)
	if _, ok := uuid.Parse(res.Header().Get(correlationIDHeader)); ok != nil {
		t.Errorf("Expected valid UUID - Got [%v]. Error: [%v]", res.Header().Get(correlationIDHeader), ok)
	}
}

func TestCUIDHeader(t *testing.T) {
	c := New(Options{
		IDType: CUID,
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	c.Handler(handler).ServeHTTP(res, req)
	expectEq(t, res.Header().Get(correlationIDHeader)[0], byte('c'))
}

func TestRandomHeader(t *testing.T) {
	c := New(Options{
		IDType: Random,
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	c.Handler(handler).ServeHTTP(res, req)
	if cID, err := strconv.Atoi(res.Header().Get(correlationIDHeader)); err != nil {
		t.Errorf("Expected random Int64 - Got [%v]. Error: [%v]", cID, err)
	}
}

func TestCustomHeader(t *testing.T) {
	c := New(Options{
		IDType:       Custom,
		CustomString: "bar",
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	c.Handler(handler).ServeHTTP(res, req)
	expectEq(t, res.Header().Get(correlationIDHeader), c.opt.CustomString)
}

func TestTimeHeader(t *testing.T) {
	c := New(Options{
		IDType: Time,
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	c.Handler(handler).ServeHTTP(res, req)
	if cID, err := strconv.Atoi(res.Header().Get(correlationIDHeader)); err != nil {
		t.Errorf("Expected Int64 containing Unix Epoch time - Got [%v]. Error: [%v]", cID, err)
	}
}

func TestHeaderForward(t *testing.T) {
	c := New(Options{
		IDType: Time,
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	uID := uuid.New().String()
	req.Header.Set(correlationIDHeader, uID)

	c.Handler(handler).ServeHTTP(res, req)
	expectEq(t, res.Header().Get(correlationIDHeader), uID)
}

/* Test Helpers */
func expectEq(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected [%v] (type %v) - Got [%v] (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func expectNeq(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Expected [%v] (type %v) - Got [%v] (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

/* Benchmarks */
func BenchmarkUUID(b *testing.B) {
	c := New(Options{})
	res := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
		c.Handler(handler).ServeHTTP(res, req)
	}
}

func BenchmarkCUID(b *testing.B) {
	c := New(Options{
		IDType: CUID,
	})
	res := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
		c.Handler(handler).ServeHTTP(res, req)
	}
}

func BenchmarkRandom(b *testing.B) {
	c := New(Options{
		IDType: Random,
	})
	res := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
		c.Handler(handler).ServeHTTP(res, req)
	}
}

func BenchmarkTime(b *testing.B) {
	c := New(Options{
		IDType: Time,
	})
	res := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
		c.Handler(handler).ServeHTTP(res, req)
	}
}

func BenchmarkCustom(b *testing.B) {
	c := New(Options{
		IDType:       Custom,
		CustomString: "foo",
	})
	res := httptest.NewRecorder()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "http://example.com/foo", nil)
		c.Handler(handler).ServeHTTP(res, req)
	}
}
