package main

import (
	"reflect"
	"testing"
)

var RE_BASECASES []string = []string{
	`$ cd /`,
	`$ cd ..`,
	`$ cd adir`,
	`$ cd anotherdir`,
	`$ ls`,
	`dir a`,
	`dir anotherdir`,
	`1234 file.txt`,
	`1000000 file`,
}

func TestRECDRegExp(t *testing.T) {
	c := map[string][][]string{
		`$ cd /`: {
			{
				`$ cd /`,
				`/`,
			},
		},
		`$ cd ..`: {
			{
				`$ cd ..`,
				`..`,
			},
		},
		`$ cd adir`: {
			{
				`$ cd adir`,
				`adir`,
			},
		},
		`$ cd anotherdir`: {
			{
				`$ cd anotherdir`,
				`anotherdir`,
			},
		},
	}

	fillCases(c)
	for input, expected := range c {
		if result := RE_CD.FindAllStringSubmatch(input, 1); !reflect.DeepEqual(result, expected) {
			t.Errorf("RegExp for CD failed for %s, expected %s but got %s",
				input, expected, result)
		}
	}
}

func TestRELSRegExp(t *testing.T) {
	c := map[string][][]string{
		`$ ls`: {
			{
				`$ ls`,
			},
		},
	}

	fillCases(c)
	for input, expected := range c {
		if result := RE_LS.FindAllStringSubmatch(input, 1); !reflect.DeepEqual(result, expected) {
			t.Errorf("RegExp for CD failed for %s, expected %s but got %s",
				input, expected, result)
		}
	}
}

func TestREDIRRegExp(t *testing.T) {
	c := map[string][][]string{
		`dir a`: {
			{
				`dir a`,
				`a`,
			},
		},
		`dir anotherdir`: {
			{
				`dir anotherdir`,
				`anotherdir`,
			},
		},
	}

	fillCases(c)
	for input, expected := range c {
		if result := RE_DIR.FindAllStringSubmatch(input, 1); !reflect.DeepEqual(result, expected) {
			t.Errorf("RegExp for CD failed for %s, expected %s but got %s",
				input, expected, result)
		}
	}
}

func TestREFILERegExp(t *testing.T) {
	c := map[string][][]string{
		`1234 file.txt`: {
			{
				`1234 file.txt`,
				`1234`,
				`file.txt`,
			},
		},
		`1000000 file`: {
			{
				`1000000 file`,
				`1000000`,
				`file`,
			},
		},
	}

	fillCases(c)
	for input, expected := range c {
		if result := RE_FILE.FindAllStringSubmatch(input, 1); !reflect.DeepEqual(result, expected) {
			t.Errorf("RegExp for CD failed for %s, expected %s but got %s",
				input, expected, result)
		}
	}
}

func fillCases(c map[string][][]string) {
	for _, tc := range RE_BASECASES {
		if _, exists := c[tc]; !exists {
			c[tc] = nil
		}
	}
}
