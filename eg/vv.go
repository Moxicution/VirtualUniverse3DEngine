// Copyright © 2018-2024 Galvanized Logic Inc.

package main

import (
	"fmt"
	"log/slog"

	"github.com/gazed/vu/internal/render/vk"
)

// vv checks that Vulkan is available, that the bindings are working
// and that it can be initialized and shutdown. The vulkan version
// is dumped to the console.
//
// CONTROLS: NA
func vv() {
	fmt.Printf("eg.vv checks that a vulkan instance can be created.\n")
	ver, err := vk.EnumerateInstanceVersion()

	if err != nil {
		slog.Error("vk.EnumerateInstanceVersion() failed", "err", err.Error())
		return
	}
	slog.Info("vulkan API", "version", vulkanVersionStr(ver))

	// create a vulkan instance
	appInfo := vk.ApplicationInfo{
		PApplicationName:   "vv",
		ApplicationVersion: vk.MAKE_VERSION(1, 0, 0),
		EngineVersion:      vk.MAKE_VERSION(1, 0, 0),
		ApiVersion:         vk.MAKE_VERSION(1, 2, 0),
	}
	ci := vk.InstanceCreateInfo{
		PApplicationInfo: &appInfo,
	}
	instance, err := vk.CreateInstance(&ci, nil)
	if err != nil {
		slog.Error("vk.CreateInstance failed", "err", err.Error())
		return
	}

	devices, err := vk.EnumeratePhysicalDevices(instance)
	if err != nil {
		slog.Error("vk.PhysicalDevice failed", "err", err.Error())
	}

	defer vk.DestroyInstance(instance, nil)
	slog.Info("vulkan started", "instance", instance)
	fmt.Println([]any{devices[0]}...)
	slog.Info("Devices")
	slog.Info("vulkan shutdown")
}

// Helper to extract parts of the Vulkan version and convert to a string
func vulkanVersionStr(version uint32) string {
	return fmt.Sprintf("%d.%d.%d",
		vk.API_VERSION_MAJOR(version),
		vk.API_VERSION_MINOR(version),
		vk.API_VERSION_PATCH(version),
	)
}
