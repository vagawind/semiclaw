package sandbox

import "testing"

func TestIsStandardTemplateRecognizesProviderScopedName(t *testing.T) {
	for _, name := range []string{"semiclaw", "team/semiclaw", "project-b89e/SemiClaw"} {
		if !isStandardTemplate(name) {
			t.Fatalf("expected %q to identify the SemiClaw standard template", name)
		}
	}
	if isStandardTemplate("semiclaw-custom") {
		t.Fatal("custom template must not be treated as the standard template")
	}
	if isStandardTemplate(DesktopTemplateName) {
		t.Fatal("the desktop sibling must not be classified as the CLI standard template")
	}
}

func TestIsDesktopTemplateRecognizesProviderScopedName(t *testing.T) {
	for _, name := range []string{"semiclaw-desktop", "team/semiclaw-desktop", "project-b89e/SemiClaw-Desktop"} {
		if !isDesktopTemplate(name) {
			t.Fatalf("expected %q to identify the SemiClaw desktop template", name)
		}
	}
	if isDesktopTemplate("semiclaw") {
		t.Fatal("the CLI template must not be classified as desktop")
	}
	if isDesktopTemplate("semiclaw-desktop-custom") {
		t.Fatal("a similarly prefixed custom name must not be the desktop template")
	}
}

func TestClassifySemiClawTemplatePrefersNameOverImage(t *testing.T) {
	standard, desktop := classifySemiClawTemplate(DesktopTemplateName, DefaultDockerImage)
	if standard || !desktop {
		t.Fatalf(
			"named desktop template must be desktop even if the image repo matches CLI, got standard=%v desktop=%v",
			standard, desktop,
		)
	}
	standard, desktop = classifySemiClawTemplate(StandardTemplateName, DefaultDesktopDockerImage)
	if !standard || desktop {
		t.Fatalf(
			"named CLI template must stay CLI even if the image tag is desktop, got standard=%v desktop=%v",
			standard, desktop,
		)
	}
}

func TestClassifySemiClawTemplateNamelessImageUsesTag(t *testing.T) {
	standard, desktop := classifySemiClawTemplate("", DefaultDockerImage)
	if !standard || desktop {
		t.Fatalf("CLI image with no name must be standard, got standard=%v desktop=%v", standard, desktop)
	}
	standard, desktop = classifySemiClawTemplate("", DefaultDesktopDockerImage)
	if standard || !desktop {
		t.Fatalf("desktop image with no name must be desktop, got standard=%v desktop=%v", standard, desktop)
	}
	standard, desktop = classifySemiClawTemplate("", DefaultCubeDesktopTemplateImage)
	if standard || !desktop {
		t.Fatalf("Cube desktop image with no name must be desktop, got standard=%v desktop=%v", standard, desktop)
	}
}
