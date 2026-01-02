package test

import (
	"testing"

	"github.com/wellingtonlope/ticket-api/test/steps"
)

func TestRegisterFeature(t *testing.T) {
	steps.RegisterFeature(t, "features/register.feature")
}
