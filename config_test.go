package twerge

import (
	"testing"
)

func TestArbitraryShadow(t *testing.T) {
	if !isArbitraryShadow("[inset_0_1px_0,inset_0_-1px_0]") {
		t.Error("isArbitraryShadow() should return true")
	}
	if !isArbitraryShadow("[0_35px_60px_-15px_rgba(0,0,0,0.3)]") {
		t.Error("isArbitraryShadow() should return true")
	}
	if !isArbitraryShadow("[inset_0_1px_0,inset_0_-1px_0]") {
		t.Error("isArbitraryShadow() should return true")
	}
	if !isArbitraryShadow("[0_0_#00f]") {
		t.Error("isArbitraryShadow() should return true")
	}
	if !isArbitraryShadow("[.5rem_0_rgba(5,5,5,5)]") {
		t.Error("isArbitraryShadow() should return true")
	}
	if !isArbitraryShadow("[-.5rem_0_#123456]") {
		t.Error("isArbitraryShadow() should return true")
	}
	if !isArbitraryShadow("[0.5rem_-0_#123456]") {
		t.Error("isArbitraryShadow() should return true")
	}
	if !isArbitraryShadow("[0.5rem_-0.005vh_#123456]") {
		t.Error("isArbitraryShadow() should return true")
	}
	if !isArbitraryShadow("[0.5rem_-0.005vh]") {
		t.Error("isArbitraryShadow() should return true")
	}

	if isArbitraryShadow("[rgba(5,5,5,5)]") {
		t.Error("isArbitraryShadow() should return false")
	}
	if isArbitraryShadow("[#00f]") {
		t.Error("isArbitraryShadow() should return false")
	}
	if isArbitraryShadow("[something-else]") {
		t.Error("isArbitraryShadow() should return false")
	}
}
