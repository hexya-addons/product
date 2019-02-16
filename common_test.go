// Copyright 2018 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"testing"

	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/tests"
	"github.com/hexya-erp/pool/h"
	"github.com/hexya-erp/pool/m"
	"github.com/hexya-erp/pool/q"
)

func TestMain(m *testing.M) {
	tests.RunTests(m, "product")
}

type productTestData struct {
	partner1    m.PartnerSet
	uomUnit     m.ProductUomSet
	uomDozen    m.ProductUomSet
	uomDunit    m.ProductUomSet
	uomWeight   m.ProductUomSet
	product0    m.ProductProductSet
	product1    m.ProductProductSet
	product2    m.ProductProductSet
	product3    m.ProductProductSet
	product4    m.ProductProductSet
	product5    m.ProductProductSet
	product6    m.ProductProductSet
	product7    m.ProductProductSet
	product71   m.ProductProductSet
	product72   m.ProductProductSet
	product8    m.ProductProductSet
	product9    m.ProductProductSet
	product10   m.ProductProductSet
	template7   m.ProductTemplateSet
	prodAtt1    m.ProductAttributeSet
	prodAttr1V1 m.ProductAttributeValueSet
	prodAttr1V2 m.ProductAttributeValueSet
}

func getProductTestData(env models.Environment) *productTestData {
	var ptd productTestData
	ptd.partner1 = h.Partner().Create(env, h.Partner().NewData().
		SetName("Julia Agrolait").
		SetEmail("julia@agrolait.example.com"))
	ptd.uomUnit = h.ProductUom().Search(env, q.ProductUom().HexyaExternalID().Equals("product_product_uom_unit"))
	ptd.uomDozen = h.ProductUom().Search(env, q.ProductUom().HexyaExternalID().Equals("product_product_uom_dozen"))
	ptd.uomDunit = h.ProductUom().Create(env, h.ProductUom().NewData().
		SetName("DeciUnit").
		SetCategory(ptd.uomUnit.Category()).
		SetFactorInv(0.1).
		SetFactor(10).
		SetUomType("smaller").
		SetRounding(0.001))
	ptd.uomWeight = h.ProductUom().Search(env, q.ProductUom().HexyaExternalID().Equals("product_product_uom_kgm"))
	ptd.product0 = h.ProductProduct().Create(env, h.ProductProduct().NewData().
		SetName("Work").
		SetType("service").
		SetUom(ptd.uomUnit).
		SetUomPo(ptd.uomUnit))
	ptd.product1 = h.ProductProduct().Create(env, h.ProductProduct().NewData().
		SetName("Courage").
		SetType("consu").
		SetUom(ptd.uomUnit).
		SetUomPo(ptd.uomDunit))
	ptd.product2 = h.ProductProduct().Create(env, h.ProductProduct().NewData().
		SetName("Wood").
		SetUom(ptd.uomUnit).
		SetUomPo(ptd.uomUnit))
	ptd.product3 = h.ProductProduct().Create(env, h.ProductProduct().NewData().
		SetName("Stone").
		SetUom(ptd.uomDozen).
		SetUomPo(ptd.uomDozen))
	ptd.product4 = h.ProductProduct().Create(env, h.ProductProduct().NewData().
		SetName("Stick").
		SetUom(ptd.uomDozen).
		SetUomPo(ptd.uomDozen))
	ptd.product5 = h.ProductProduct().Create(env, h.ProductProduct().NewData().
		SetName("Stone Tools").
		SetUom(ptd.uomUnit).
		SetUomPo(ptd.uomUnit))
	ptd.product6 = h.ProductProduct().Create(env, h.ProductProduct().NewData().
		SetName("Door").
		SetUom(ptd.uomUnit).
		SetUomPo(ptd.uomUnit))
	ptd.prodAtt1 = h.ProductAttribute().Create(env, h.ProductAttribute().NewData().
		SetName("Color"))
	ptd.prodAttr1V1 = h.ProductAttributeValue().Create(env, h.ProductAttributeValue().NewData().
		SetName("Red").
		SetAttribute(ptd.prodAtt1))
	ptd.prodAttr1V2 = h.ProductAttributeValue().Create(env, h.ProductAttributeValue().NewData().
		SetName("Blue").
		SetAttribute(ptd.prodAtt1))
	ptd.template7 = h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
		SetName("Sofa").
		SetUom(ptd.uomUnit).
		SetUomPo(ptd.uomUnit).
		SetAttributeLines(h.ProductAttributeLine().Create(env, h.ProductAttributeLine().NewData().
			SetAttribute(ptd.prodAtt1))))
	ptd.product7 = h.ProductProduct().Create(env, h.ProductProduct().NewData().
		SetProductTmpl(ptd.template7))
	ptd.product71 = h.ProductProduct().Create(env, h.ProductProduct().NewData().
		SetProductTmpl(ptd.template7).
		SetAttributeValues(ptd.prodAttr1V1))
	ptd.product72 = h.ProductProduct().Create(env, h.ProductProduct().NewData().
		SetProductTmpl(ptd.template7).
		SetAttributeValues(ptd.prodAttr1V2))
	ptd.product8 = h.ProductProduct().Create(env, h.ProductProduct().NewData().
		SetName("House").
		SetUom(ptd.uomUnit).
		SetUomPo(ptd.uomUnit))
	ptd.product9 = h.ProductProduct().Create(env, h.ProductProduct().NewData().
		SetName("Paper").
		SetUom(ptd.uomUnit).
		SetUomPo(ptd.uomUnit))
	ptd.product10 = h.ProductProduct().Create(env, h.ProductProduct().NewData().
		SetName("Stone").
		SetUom(ptd.uomUnit).
		SetUomPo(ptd.uomUnit))
	return &ptd
}
