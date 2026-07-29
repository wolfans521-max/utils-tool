package card

import (
	"encoding/json"
	"testing"
)

// TestNewItemCellSpan 校验工厂方法创建的 CellSpan 基础字段
func TestNewItemCellSpan(t *testing.T) {
	cell := NewItemCellSpan()

	if cell == nil {
		t.Fatalf("NewItemCellSpan returned nil")
	}

	// 匿名嵌套的 ItemBase 应该设置 category 为 cell
	if cell.Category != CategoryCell {
		t.Fatalf("Category expected %q, got %q", CategoryCell, cell.Category)
	}

	// cellBase.Type 应为 span
	if cell.Type != "span" {
		t.Fatalf("Type expected span, got %s", cell.Type)
	}
}

// TestItemCellSpanJSONStructure 校验 CellSpan 的 JSON 结构
func TestItemCellSpanJSONStructure(t *testing.T) {
	cell := NewItemCellSpan()

	b, err := json.Marshal(cell)
	if err != nil {
		t.Fatalf("marshal CellSpan failed: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal to map failed: %v", err)
	}

	// 顶层 category / type
	if v, ok := m["category"]; !ok || v != string(CategoryCell) {
		t.Fatalf("category expect %q, got %v (ok=%v)", CategoryCell, v, ok)
	}
	if v, ok := m["type"]; !ok || v != "span" {
		t.Fatalf("type expect span, got %v (ok=%v)", v, ok)
	}
}
