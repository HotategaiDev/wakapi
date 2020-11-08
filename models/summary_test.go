package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSummary_FillUnknown(t *testing.T) {
	testDuration := 10 * time.Minute

	sut := &Summary{
		Projects: []*SummaryItem{
			{
				Type: SummaryProject,
				Key:  "wakapi",
				// hack to work around the issue that the total time of a summary item is mistakenly represented in seconds
				Total: testDuration / time.Second,
			},
		},
	}

	sut.FillUnknown()

	itemLists := [][]*SummaryItem{
		sut.Machines,
		sut.OperatingSystems,
		sut.Languages,
		sut.Editors,
	}
	for _, l := range itemLists {
		assert.Len(t, l, 1)
		assert.Equal(t, UnknownSummaryKey, l[0].Key)
		assert.Equal(t, testDuration, l[0].Total)
	}
}

func TestSummary_TotalTimeBy(t *testing.T) {
	testDuration1, testDuration2, testDuration3 := 10*time.Minute, 5*time.Minute, 20*time.Minute

	sut := &Summary{
		Projects: []*SummaryItem{
			{
				Type: SummaryProject,
				Key:  "wakapi",
				// hack to work around the issue that the total time of a summary item is mistakenly represented in seconds
				Total: testDuration1 / time.Second,
			},
			{
				Type:  SummaryProject,
				Key:   "anchr",
				Total: testDuration2 / time.Second,
			},
		},
		Languages: []*SummaryItem{
			{
				Type:  SummaryLanguage,
				Key:   "Go",
				Total: testDuration3 / time.Second,
			},
		},
	}

	assert.Equal(t, testDuration1+testDuration2, sut.TotalTimeBy(SummaryProject))
	assert.Equal(t, testDuration3, sut.TotalTimeBy(SummaryLanguage))
	assert.Zero(t, sut.TotalTimeBy(SummaryEditor))
	assert.Zero(t, sut.TotalTimeBy(SummaryMachine))
	assert.Zero(t, sut.TotalTimeBy(SummaryOS))
}

func TestSummary_TotalTimeByFilters(t *testing.T) {
	testDuration1, testDuration2, testDuration3 := 10*time.Minute, 5*time.Minute, 20*time.Minute

	sut := &Summary{
		Projects: []*SummaryItem{
			{
				Type: SummaryProject,
				Key:  "wakapi",
				// hack to work around the issue that the total time of a summary item is mistakenly represented in seconds
				Total: testDuration1 / time.Second,
			},
			{
				Type:  SummaryProject,
				Key:   "anchr",
				Total: testDuration2 / time.Second,
			},
		},
		Languages: []*SummaryItem{
			{
				Type:  SummaryLanguage,
				Key:   "Go",
				Total: testDuration3 / time.Second,
			},
		},
	}

	filters1 := &Filters{Project: "wakapi"}
	filters2 := &Filters{Project: "wakapi", Language: "Go"} // filters have OR logic
	filters3 := &Filters{}

	assert.Equal(t, testDuration1, sut.TotalTimeByFilters(filters1))
	assert.Equal(t, testDuration1+testDuration3, sut.TotalTimeByFilters(filters2))
	assert.Zero(t, sut.TotalTimeByFilters(filters3))
}

func TestSummary_WithResolvedAliases(t *testing.T) {
	testDuration1, testDuration2, testDuration3, testDuration4 := 10*time.Minute, 5*time.Minute, 1*time.Minute, 20*time.Minute

	var resolver AliasResolver = func(t uint8, k string) string {
		switch t {
		case SummaryProject:
			switch k {
			case "wakapi-mobile":
				return "wakapi"
			}
		case SummaryLanguage:
			switch k {
			case "Java 8":
				return "Java"
			}
		}
		return k
	}

	sut := &Summary{
		Projects: []*SummaryItem{
			{
				Type:  SummaryProject,
				Key:   "wakapi",
				Total: testDuration1 / time.Second,
			},
			{
				Type:  SummaryProject,
				Key:   "wakapi-mobile",
				Total: testDuration2 / time.Second,
			},
			{
				Type:  SummaryProject,
				Key:   "anchr",
				Total: testDuration3 / time.Second,
			},
		},
		Languages: []*SummaryItem{
			{
				Type:  SummaryLanguage,
				Key:   "Java 8",
				Total: testDuration4 / time.Second,
			},
		},
	}

	sut = sut.WithResolvedAliases(resolver)

	assert.Equal(t, testDuration1+testDuration2, sut.TotalTimeByKey(SummaryProject, "wakapi"))
	assert.Zero(t, sut.TotalTimeByKey(SummaryProject, "wakapi-mobile"))
	assert.Equal(t, testDuration3, sut.TotalTimeByKey(SummaryProject, "anchr"))
	assert.Equal(t, testDuration4, sut.TotalTimeByKey(SummaryLanguage, "Java"))
	assert.Zero(t, sut.TotalTimeByKey(SummaryLanguage, "wakapi"))
	assert.Zero(t, sut.TotalTimeByKey(SummaryProject, "Java 8"))
	assert.Len(t, sut.Projects, 2)
	assert.Len(t, sut.Languages, 1)
	assert.Empty(t, sut.Editors)
	assert.Empty(t, sut.OperatingSystems)
	assert.Empty(t, sut.Machines)
}
