const std = @import("std");
const builtin = @import("builtin");

const LIB_NAME = "mymath";

const goosConv = std.StaticStringMap([]const u8).initComptime(.{
    .{ "darwin", "macos" },
    .{ "windows", "windows" },
    .{ "linux", "linux" },
});

const goarchConv = std.StaticStringMap([]const u8).initComptime(.{
    .{ "amd64", "x86_64" },
    .{ "arm64", "aarch64" },
    .{ "arm", "arm" },
    .{ "386", "i386" },
});

var sharedLib: []const u8 = undefined;
var zigLib: []const u8 = undefined;
var cwd: std.Io.Dir = undefined;
var io: std.Io = undefined;

fn copyAndRenameFiles(step: *std.Build.Step, options: std.Build.Step.MakeOptions) !void {
    _ = step;
    _ = options;
    try cwd.copyFile(sharedLib, cwd, zigLib, io, .{});
}

pub fn build(b: *std.Build) !void {
    const allocator = b.allocator;
    io = b.graph.io;

    // Default OS and ARCH
    var os: []const u8 = @tagName(builtin.os.tag);
    var arch: []const u8 = @tagName(builtin.target.cpu.arch);

    // Check environment variables for override
    const env_map = b.graph.environ_map;

    if (env_map.get("GOOS")) |val| {
        if (goosConv.get(val)) |v| {
            os = v;
        }
    }

    if (env_map.get("GOARCH")) |val| {
        if (goarchConv.get(val)) |v| {
            arch = v;
        }
    }

    const target = try std.fmt.allocPrint(allocator, "{s}-{s}", .{ arch, os });
    defer allocator.free(target);
    const query = std.Target.Query.parse(.{
        .arch_os_abi = target,
    }) catch @panic("Invalid target string");

    const lib = b.addLibrary(.{
        .name = LIB_NAME,
        .linkage = .dynamic,
        .root_module = b.createModule(.{
            .root_source_file = b.path("src/math.zig"),
            .target = b.resolveTargetQuery(query),
            .optimize = b.standardOptimizeOption(.{}),
        }),
    });

    // This makes the output file available in zig-out/(lib|bin)
    const install_step = b.addInstallArtifact(lib, .{});

    // Find the shared library
    if (std.mem.eql(u8, os, "macos")) {
        sharedLib = try std.fmt.allocPrint(allocator, "{s}/lib/lib{s}.dylib", .{ b.install_path, LIB_NAME });
    } else if (std.mem.eql(u8, os, "windows")) {
        sharedLib = try std.fmt.allocPrint(allocator, "{s}/bin/{s}.dll", .{ b.install_path, LIB_NAME });
    } else if (std.mem.eql(u8, os, "linux")) {
        sharedLib = try std.fmt.allocPrint(allocator, "{s}/lib/lib{s}.so", .{ b.install_path, LIB_NAME });
    }

//     std.debug.print("Shared file: {s}\n", .{sharedLib});

    // Gzip up the file and rename it to <name>.shared (Gzip is currently removed from zig std lib)

    // Copy shared library to install_path and rename extension to .shared
    zigLib = try std.fmt.allocPrint(allocator, "{s}/{s}.shared", .{ b.install_path, LIB_NAME });
//     std.debug.print("Final file: {s}\n", .{zigLib});

    cwd = std.Io.Dir.cwd();

    const copy_step = b.step("go-build", "Build Library for Go embed");
    copy_step.makeFn = copyAndRenameFiles;
    copy_step.dependOn(&install_step.step);
}
