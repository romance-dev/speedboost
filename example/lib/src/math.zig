const std = @import("std");

export fn add(a: i64, b: i64) callconv(.c) i64 {
    return a + b;
}

export fn multiply(a: f64, b: f64) callconv(.c) f64 {
    return a * b;
}

// Use 'export' to make it visible to the linker.
// 'callconv(.C)' ensures it follows the target's C calling convention.
export fn add_numbers(a: i32, b: i32) callconv(.c) i32 {
    return a + b;
}

// For structs, use 'extern struct' to guarantee C-compatible memory layout.
pub const MyData = extern struct {
    id: u32,
    value: f32,
};

export fn process_data(data: *const MyData) callconv(.c) f32 {
    return @as(f32, @floatFromInt(data.id)) + data.value;
}
