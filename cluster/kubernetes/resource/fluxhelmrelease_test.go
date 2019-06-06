package resource

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/weaveworks/flux/resource"
)

func TestSortedContainers(t *testing.T) {
	unordered, expected := map[string]imageAndSetter{
		"ZZZ": 					{},
		"AAA": 					{},
		"FFF": 					{},
		ReleaseContainerName:	{},
	}, []string{ReleaseContainerName, "AAA", "FFF", "ZZZ"}

	actual := sorted_containers(unordered)

	assert.Equal(t, expected, actual)
}

func TestParseImageOnlyFormat(t *testing.T) {
	expectedImage := "bitnami/mariadb:10.1.30-r1"
	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    first: post
    image: ` + expectedImage + `
    persistence:
      enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Errorf("expected 1 container; got %#v", containers)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}
}

func TestParseImageTagFormat(t *testing.T) {
	expectedImageName := "bitnami/mariadb"
	expectedImageTag := "10.1.30-r1"
	expectedImage := expectedImageName + ":" + expectedImageTag

	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    first: post
    image: ` + expectedImageName + `
    tag: ` + expectedImageTag + `
    persistence:
      enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Errorf("expected 1 container; got %#v", containers)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}
}

func TestParseNamedImageFormat(t *testing.T) {
	expectedContainer := "db"
	expectedImage := "bitnami/mariadb:10.1.30-r1"
	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    ` + expectedContainer + `:
      first: post
      image: ` + expectedImage + `
      persistence:
        enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}

	newImage := containers[0].Image.WithNewTag("some-other-tag")
	if err := fhr.SetContainerImage(expectedContainer, newImage); err != nil {
		t.Error(err)
	}

	containers = fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image = containers[0].Image.String()
	if image != newImage.String() {
		t.Errorf("expected container image %q, got %q", newImage.String(), image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}
}

func TestParseNamedImageTagFormat(t *testing.T) {
	expectedContainer := "db"
	expectedImageName := "bitnami/mariadb"
	expectedImageTag := "10.1.30-r1"
	expectedImage := expectedImageName + ":" + expectedImageTag

	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    other:
      not: "containing image"
    ` + expectedContainer + `:
      first: post
      image: ` + expectedImageName + `
      tag: ` + expectedImageTag + `
      persistence:
        enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}

	newImage := containers[0].Image.WithNewTag("some-other-tag")
	if err := fhr.SetContainerImage(expectedContainer, newImage); err != nil {
		t.Error(err)
	}

	containers = fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image = containers[0].Image.String()
	if image != newImage.String() {
		t.Errorf("expected container image %q, got %q", newImage.String(), image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}
}

func TestParseImageObjectFormat(t *testing.T) {
	expectedImageName := "bitnami/mariadb"
	expectedImageTag := "10.1.30-r1"
	expectedImage := expectedImageName + ":" + expectedImageTag

	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    first: post
    image:
      repository: ` + expectedImageName + `
      tag: ` + expectedImageTag + `
    persistence:
      enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Errorf("expected 1 container; got %#v", containers)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}
}

func TestParseNamedImageObjectFormat(t *testing.T) {
	expectedContainer := "db"
	expectedImageName := "bitnami/mariadb"
	expectedImageTag := "10.1.30-r1"
	expectedImage := expectedImageName + ":" + expectedImageTag

	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    other:
      not: "containing image"
    ` + expectedContainer + `:
      first: post
      image:
        repository: ` + expectedImageName + `
        tag: ` + expectedImageTag + `
      persistence:
        enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}

	newImage := containers[0].Image.WithNewTag("some-other-tag")
	if err := fhr.SetContainerImage(expectedContainer, newImage); err != nil {
		t.Error(err)
	}

	containers = fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image = containers[0].Image.String()
	if image != newImage.String() {
		t.Errorf("expected container image %q, got %q", newImage.String(), image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}
}

func TestParseMappedImageOnly(t *testing.T) {
	expectedContainer := "mariadb"
	expectedImage := "bitnami/mariadb:10.1.30-r1"
	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  annotations:
    ` + ImageRepositoryPrefix + expectedContainer + `: customRepository
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    first: post
    customRepository: ` + expectedImage + `
    persistence:
      enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Errorf("expected 1 container; got %#v", containers)
	}
	container := containers[0].Name
	if container != expectedContainer {
		t.Errorf("expected container container %q, got %q", expectedContainer, container)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}

	newImage := containers[0].Image.WithNewTag("some-other-tag")
	if err := fhr.SetContainerImage(expectedContainer, newImage); err != nil {
		t.Error(err)
	}

	containers = fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image = containers[0].Image.String()
	if image != newImage.String() {
		t.Errorf("expected container image %q, got %q", newImage.String(), image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}
}

func TestParseMappedImageTag(t *testing.T) {
	expectedContainer := "mariadb"
	expectedImageName := "bitnami/mariadb"
	expectedImageTag := "10.1.30-r1"
	expectedImage := expectedImageName + ":" + expectedImageTag
	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  annotations:
    ` + ImageRepositoryPrefix + expectedContainer + `: customRepository
    ` + ImageTagPrefix + expectedContainer + `: customTag
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    first: post
    customRepository: ` + expectedImageName + `
    customTag: ` + expectedImageTag + `
    persistence:
      enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 1 {
		t.Errorf("expected 1 container; got %#v", containers)
	}
	container := containers[0].Name
	if container != expectedContainer {
		t.Errorf("expected container container %q, got %q", expectedContainer, container)
	}
	image := containers[0].Image.String()
	if image != expectedImage {
		t.Errorf("expected container image %q, got %q", expectedImage, image)
	}

	newImage := containers[0].Image.WithNewTag("some-other-tag")
	if err := fhr.SetContainerImage(expectedContainer, newImage); err != nil {
		t.Error(err)
	}

	containers = fhr.Containers()
	if len(containers) != 1 {
		t.Fatalf("expected 1 container; got %#v", containers)
	}
	image = containers[0].Image.String()
	if image != newImage.String() {
		t.Errorf("expected container image %q, got %q", newImage.String(), image)
	}
	if containers[0].Name != expectedContainer {
		t.Errorf("expected container name %q, got %q", expectedContainer, containers[0].Name)
	}
}

func TestParseMappedTagOnly(t *testing.T) {
	container := "mariadb"
	imageTag := "10.1.30-r1"
	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: mariadb
  namespace: maria
  annotations:
    ` + ImageTagPrefix + container + `: customTag
  labels:
    chart: mariadb
spec:
  chartGitPath: mariadb
  values:
    first: post
    customTag: ` + imageTag + `
    persistence:
      enabled: false
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["maria:fluxhelmrelease/mariadb"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != 0 {
		t.Errorf("expected 0 container; got %#v", containers)
	}
}

func TestParseMappedImageFormats(t *testing.T) {

	expected := []struct {
		name, image, tag string
	}{
		{"AAA", "repo/imageOne", "tagOne"},
		{"CCC", "repo/imageTwo", "tagTwo"},
		{"ZZZ", "repo/imageThree", "tagThree"},
	}

	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: test
  namespace: test
  annotations:
    # Top level mapping
    ` + ImageRepositoryPrefix + expected[0].name + `: customRepository
    ` + ImageTagPrefix + expected[0].name + `: customTag
    # Sub level mapping
    ` + ImageRepositoryPrefix + expected[1].name + `: ` + expected[1].name + `.customRepository
    ` + ImageTagPrefix + expected[1].name + `: ` + expected[1].name + `.customTag
    # Sub level mapping 2
    ` + ImageRepositoryPrefix + expected[2].name + `: ` + expected[2].name + `.image.customRepository
spec:
  chartGitPath: test
  values:
    # Top level image
    customRepository: ` + expected[0].image + `
    customTag: ` + expected[0].tag + `

    # Sub level image
    ` + expected[1].name + `:
      customRepository: ` + expected[1].image + `
      customTag: ` + expected[1].tag + `

    # Sub level image 2
    ` + expected[2].name + `:
      image:
        customRepository: ` + expected[2].image + `:` + expected[2].tag + `
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["test:fluxhelmrelease/test"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != len(expected) {
		t.Fatalf("expected %d containers, got %d, %#v", len(expected), len(containers), containers)
	}
	for i, c0 := range expected {
		c1 := containers[i]
		if c1.Name != c0.name {
			t.Errorf("names do not match %q != %q", c0, c1)
		}
		c0image := fmt.Sprintf("%s:%s", c0.image, c0.tag)
		if c1.Image.String() != c0image {
			t.Errorf("images do not match %q != %q", c0image, c1.Image.String())
		}
	}
}

func TestParseAllFormatsInOne(t *testing.T) {

	type container struct {
		name, image, tag string
	}

	// *NB* the containers will be calculated based on the order
	//  1. the entry for 'image' if present
	//  2. the order of the keys in `values`.
	//
	// To avoid having to mess around later, I have cooked the order
	// of these so they can be compared directly to the return value.
	expected := []container{
		{ReleaseContainerName, "repo/imageOne", "tagOne"},
		{"AAA", "repo/imageTwo", "tagTwo"},
		{"KKK", "repo/imageThree", "tagThree"},
		{"ZZZ", "repo/imageFour", "tagFour"},
	}

	doc := `---
apiVersion: helm.integrations.flux.weave.works/v1alpha2
kind: FluxHelmRelease
metadata:
  name: test
  namespace: test
  annotations:
    ` + ImageRepositoryPrefix + expected[3].name + `: ` + expected[3].name +`.customRepository
    ` + ImageTagPrefix + expected[3].name + `: ` + expected[3].name +`.customTag
spec:
  chartGitPath: test
  values:
    # top-level image
    image: ` + expected[0].image + ":" + expected[0].tag + `

    # under .container, as image and tag entries
    ` + expected[1].name + `:
      image: ` + expected[1].image + `
      tag: ` + expected[1].tag + `

    # under .container.image, as repository and tag entries
    ` + expected[2].name + `:
      image:
        repository: ` + expected[2].image + `
        tag: ` + expected[2].tag + `
      persistence:
        enabled: false

    # mapped by user via annotations
    ` + expected[3].name + `:
      customRepository: ` + expected[3].image + `
      customTag: ` + expected[3].tag + `
`

	resources, err := ParseMultidoc([]byte(doc), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, ok := resources["test:fluxhelmrelease/test"]
	if !ok {
		t.Fatalf("expected resource not found; instead got %#v", resources)
	}
	fhr, ok := res.(resource.Workload)
	if !ok {
		t.Fatalf("expected resource to be a Workload, instead got %#v", res)
	}

	containers := fhr.Containers()
	if len(containers) != len(expected) {
		t.Fatalf("expected %d containers, got %d", len(expected), len(containers))
	}
	for i, c0 := range expected {
		c1 := containers[i]
		if c1.Name != c0.name {
			t.Errorf("names do not match %q != %q", c0, c1)
		}
		c0image := fmt.Sprintf("%s:%s", c0.image, c0.tag)
		if c1.Image.String() != c0image {
			t.Errorf("images do not match %q != %q", c0image, c1.Image.String())
		}
	}
}
