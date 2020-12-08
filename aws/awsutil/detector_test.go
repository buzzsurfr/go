package awsutil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test slice should be in the same order as the ComputeType consts
var detectOutputTests = []DetectOutput{
	{Type: ComputeUnknown},
	{Type: ComputeEC2},
	{Type: ComputeDocker},
	{Type: ComputeKubernetes},
	{Type: ComputeECS},
	{Type: ComputeEKS},
	{Type: ComputeLambda},
}

func TestIsEC2(t *testing.T) {
	for i, dot := range detectOutputTests {
		if ComputeType(i) == ComputeEC2 {
			testname := fmt.Sprintf("IsEC2 is true for %s", dot)
			t.Run(testname, func(t *testing.T) {
				assert.True(t, IsEC2(dot))
			})
		} else {
			testname := fmt.Sprintf("IsEC2 is false for %s", dot)
			t.Run(testname, func(t *testing.T) {
				assert.False(t, IsEC2(dot))
			})
		}
	}
}

func TestIsDocker(t *testing.T) {
	for i, dot := range detectOutputTests {
		if ComputeType(i) == ComputeDocker {
			testname := fmt.Sprintf("IsDocker is true for %s", dot)
			t.Run(testname, func(t *testing.T) {
				assert.True(t, IsDocker(dot))
			})
		} else {
			testname := fmt.Sprintf("IsDocker is false for %s", dot)
			t.Run(testname, func(t *testing.T) {
				assert.False(t, IsDocker(dot))
			})
		}
	}
}

func TestIsKubernetes(t *testing.T) {
	for i, dot := range detectOutputTests {
		if ComputeType(i) == ComputeKubernetes {
			testname := fmt.Sprintf("IsKubernetes is true for %s", dot)
			t.Run(testname, func(t *testing.T) {
				assert.True(t, IsKubernetes(dot))
			})
		} else {
			testname := fmt.Sprintf("IsKubernetes is false for %s", dot)
			t.Run(testname, func(t *testing.T) {
				assert.False(t, IsKubernetes(dot))
			})
		}
	}
}

func TestIsECS(t *testing.T) {
	for i, dot := range detectOutputTests {
		if ComputeType(i) == ComputeECS {
			testname := fmt.Sprintf("IsECS is true for %s", dot)
			t.Run(testname, func(t *testing.T) {
				assert.True(t, IsECS(dot))
			})
		} else {
			testname := fmt.Sprintf("IsECS is false for %s", dot)
			t.Run(testname, func(t *testing.T) {
				assert.False(t, IsECS(dot))
			})
		}
	}
}

func TestIsEKS(t *testing.T) {
	for i, dot := range detectOutputTests {
		if ComputeType(i) == ComputeEKS {
			testname := fmt.Sprintf("IsEKS is true for %s", dot)
			t.Run(testname, func(t *testing.T) {
				assert.True(t, IsEKS(dot))
			})
		} else {
			testname := fmt.Sprintf("IsEKS is false for %s", dot)
			t.Run(testname, func(t *testing.T) {
				assert.False(t, IsEKS(dot))
			})
		}
	}
}

func TestIsLambda(t *testing.T) {
	for i, dot := range detectOutputTests {
		if ComputeType(i) == ComputeLambda {
			testname := fmt.Sprintf("IsLambda is true for %s", dot)
			t.Run(testname, func(t *testing.T) {
				assert.True(t, IsLambda(dot))
			})
		} else {
			testname := fmt.Sprintf("IsLambda is false for %s", dot)
			t.Run(testname, func(t *testing.T) {
				assert.False(t, IsLambda(dot))
			})
		}
	}
}
