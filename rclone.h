// if you change this file, you MUST run "go clean -cache -testcache"
// before rebuilding.
// See: https://github.com/golang/go/issues/24355

// call go to rust
extern void rust_insert_file_from_go(const char *json);
extern void rust_insert_files_from_go(const char *json);
extern void rust_insert_filecache_from_go(const char *json,
                                          const unsigned char *buf,
                                          long long len);

// call rust to go is handled by the extern "C" block in lib.rs
// which is copied from the .h file that comes out of the Go build
// pub fn GoListJSON(path: *const c_char);
// pub fn GoFetchFiledata(path: *const c_char, startbytepos: c_longlong,
// endbytepos: c_longlong);