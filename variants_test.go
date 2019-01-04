// Copyright 2018 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"testing"

	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/models/operator"
	"github.com/hexya-erp/hexya/src/models/security"
	"github.com/hexya-erp/pool/h"
	"github.com/hexya-erp/pool/q"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVariantsSearch(t *testing.T) {
	Convey("Testing search on variants", t, func() {
		So(models.SimulateInNewEnvironment(security.SuperUserID, func(env models.Environment) {
			sizeAttr := h.ProductAttribute().Create(env, h.ProductAttribute().NewData().SetName("Size"))
			h.ProductAttributeValue().Create(env, h.ProductAttributeValue().NewData().
				SetName("S").
				SetAttribute(sizeAttr))
			h.ProductAttributeValue().Create(env, h.ProductAttributeValue().NewData().
				SetName("M").
				SetAttribute(sizeAttr))
			sizeAttreValueL := h.ProductAttributeValue().Create(env, h.ProductAttributeValue().NewData().
				SetName("L").
				SetAttribute(sizeAttr))
			productShirtTemplate :=
				h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Shirt").
					SetAttributeLines(h.ProductAttributeLine().Create(env, h.ProductAttributeLine().NewData().
						SetAttribute(sizeAttr).
						SetValues(sizeAttreValueL))))
			Convey("Test Attribute line search", func() {
				searchNotToBeFound := h.ProductTemplate().Search(env, q.ProductTemplate().AttributeLines().ContainsEval("M"))
				So(productShirtTemplate.Intersect(searchNotToBeFound).IsEmpty(), ShouldBeTrue)

				searchAttribute := h.ProductTemplate().Search(env, q.ProductTemplate().AttributeLines().ContainsEval("Size"))
				So(productShirtTemplate.Intersect(searchAttribute).IsEmpty(), ShouldBeFalse)

				searchValue := h.ProductTemplate().Search(env, q.ProductTemplate().AttributeLines().ContainsEval("L"))
				So(productShirtTemplate.Intersect(searchValue).IsEmpty(), ShouldBeFalse)
			})
			Convey("Test Name Search", func() {
				productSlipTemplate := h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Slip"))
				res := h.ProductProduct().NewSet(env).SearchByName("Shirt", operator.NotIContains, q.ProductProductCondition{}, 0)
				So(res.Intersect(productSlipTemplate.ProductVariant()).IsEmpty(), ShouldBeFalse)
			})
		}), ShouldBeNil)
	})
}

func TestVariants(t *testing.T) {
	Convey("Testing variants", t, func() {
		So(models.ExecuteInNewEnvironment(security.SuperUserID, func(env models.Environment) {
			sizeAttr := h.ProductAttribute().Create(env, h.ProductAttribute().NewData().SetName("Size"))
			sizeAttreValueS := h.ProductAttributeValue().Create(env, h.ProductAttributeValue().NewData().
				SetName("S").
				SetAttribute(sizeAttr))
			sizeAttreValueM := h.ProductAttributeValue().Create(env, h.ProductAttributeValue().NewData().
				SetName("M").
				SetAttribute(sizeAttr))
			sizeAttreValueL := h.ProductAttributeValue().Create(env, h.ProductAttributeValue().NewData().
				SetName("L").
				SetAttribute(sizeAttr))
			ptd := getProductTestData(env)
			Convey("One variant, because mono value", func() {
				testTemplate := h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Sofa").
					SetUom(ptd.uomUnit).
					SetUomPo(ptd.uomUnit).
					SetAttributeLines(h.ProductAttributeLine().Create(env,
						h.ProductAttributeLine().NewData().
							SetAttribute(sizeAttr).
							SetValues(sizeAttreValueS))))
				So(testTemplate.ProductVariants().Len(), ShouldEqual, 1)
				So(testTemplate.ProductVariant().AttributeValues().Len(), ShouldEqual, 1)
				So(testTemplate.ProductVariant().AttributeValues().Equals(sizeAttreValueS), ShouldBeTrue)
			})
			Convey("One variant, because only 1 combination is possible", func() {
				testTemplate := h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Sofa").
					SetUom(ptd.uomUnit).
					SetUomPo(ptd.uomUnit).
					SetAttributeLines(h.ProductAttributeLine().Create(env,
						h.ProductAttributeLine().NewData().
							SetAttribute(sizeAttr).
							SetValues(sizeAttreValueS)).
						Union(
							h.ProductAttributeLine().Create(env, h.ProductAttributeLine().NewData().
								SetAttribute(ptd.prodAtt1).
								SetValues(ptd.prodAttr1V2)))))
				So(testTemplate.ProductVariants().Len(), ShouldEqual, 1)
				So(testTemplate.ProductVariant().AttributeValues().Len(), ShouldEqual, 2)
				So(testTemplate.ProductVariant().AttributeValues().Equals(sizeAttreValueS.Union(ptd.prodAttr1V2)), ShouldBeTrue)
			})
			Convey("Two variants, simple matrix", func() {
				testTemplate := h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Sofa").
					SetUom(ptd.uomUnit).
					SetUomPo(ptd.uomUnit).
					SetAttributeLines(h.ProductAttributeLine().Create(env,
						h.ProductAttributeLine().NewData().
							SetAttribute(sizeAttr).
							SetValues(sizeAttreValueS.Union(sizeAttreValueM))).
						Union(
							h.ProductAttributeLine().Create(env, h.ProductAttributeLine().NewData().
								SetAttribute(ptd.prodAtt1).
								SetValues(ptd.prodAttr1V2)))))
				So(testTemplate.ProductVariants().Len(), ShouldEqual, 2)
				productVariants := h.ProductProduct().Search(env,
					q.ProductProduct().ProductTmpl().Equals(testTemplate).
						And().AttributeValues().Equals(ptd.prodAttr1V2))

				products := productVariants.Filtered(func(rs h.ProductProductSet) bool {
					return !rs.AttributeValues().Intersect(sizeAttreValueS).IsEmpty()
				})
				So(products.Len(), ShouldEqual, 1)
				So(products.AttributeValues().Len(), ShouldEqual, 2)
				So(products.AttributeValues().Equals(sizeAttreValueS.Union(ptd.prodAttr1V2)), ShouldBeTrue)

				products = productVariants.Filtered(func(rs h.ProductProductSet) bool {
					return !rs.AttributeValues().Intersect(sizeAttreValueM).IsEmpty()
				})
				So(products.Len(), ShouldEqual, 1)
				So(products.AttributeValues().Len(), ShouldEqual, 2)
				So(products.AttributeValues().Equals(sizeAttreValueM.Union(ptd.prodAttr1V2)), ShouldBeTrue)
			})
			Convey("Value matrix: 2x3 values", func() {
				testTemplate := h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Sofa").
					SetUom(ptd.uomUnit).
					SetUomPo(ptd.uomUnit).
					SetAttributeLines(h.ProductAttributeLine().Create(env,
						h.ProductAttributeLine().NewData().
							SetAttribute(sizeAttr).
							SetValues(sizeAttreValueS.Union(sizeAttreValueM).Union(sizeAttreValueL))).
						Union(
							h.ProductAttributeLine().Create(env, h.ProductAttributeLine().NewData().
								SetAttribute(ptd.prodAtt1).
								SetValues(ptd.prodAttr1V1.Union(ptd.prodAttr1V2))))))
				So(testTemplate.ProductVariants().Len(), ShouldEqual, 6)
				for _, value1 := range []h.ProductAttributeValueSet{ptd.prodAttr1V1, ptd.prodAttr1V2} {
					productVariants := h.ProductProduct().Search(env,
						q.ProductProduct().ProductTmpl().Equals(testTemplate).
							And().AttributeValues().Equals(value1))
					for _, value2 := range []h.ProductAttributeValueSet{sizeAttreValueS, sizeAttreValueM, sizeAttreValueL} {
						products := productVariants.Filtered(func(rs h.ProductProductSet) bool {
							return !rs.AttributeValues().Intersect(value2).IsEmpty()
						})
						So(products.Len(), ShouldEqual, 1)
						So(products.AttributeValues().Equals(value1.Union(value2)), ShouldBeTrue)
					}
				}
			})
			Convey("Creation and multi-updates", func() {
				testTemplate := h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Sofa").
					SetUom(ptd.uomUnit).
					SetUomPo(ptd.uomUnit).
					SetAttributeLines(h.ProductAttributeLine().Create(env,
						h.ProductAttributeLine().NewData().
							SetAttribute(sizeAttr).
							SetValues(sizeAttreValueS.Union(sizeAttreValueM))).
						Union(h.ProductAttributeLine().Create(env, h.ProductAttributeLine().NewData().
							SetAttribute(ptd.prodAtt1).
							SetValues(ptd.prodAttr1V1.Union(ptd.prodAttr1V2))))))
				So(testTemplate.ProductVariants().Len(), ShouldEqual, 4)
				sizeAttributeLine := testTemplate.AttributeLines().Filtered(func(rs h.ProductAttributeLineSet) bool {
					return rs.Attribute().Equals(sizeAttr)
				})
				sizeAttributeLine.SetValues(sizeAttributeLine.Values().Union(sizeAttreValueL))
				// Trigger variant updates:
				testTemplate.SetAttributeLines(testTemplate.AttributeLines())
				So(testTemplate.ProductVariants().Len(), ShouldEqual, 6)
				for _, value1 := range []h.ProductAttributeValueSet{ptd.prodAttr1V1, ptd.prodAttr1V2} {
					productVariants := h.ProductProduct().Search(env,
						q.ProductProduct().ProductTmpl().Equals(testTemplate).
							And().AttributeValues().Equals(value1))
					for _, value2 := range []h.ProductAttributeValueSet{sizeAttreValueS, sizeAttreValueM, sizeAttreValueL} {
						products := productVariants.Filtered(func(rs h.ProductProductSet) bool {
							return !rs.AttributeValues().Intersect(value2).IsEmpty()
						})
						So(products.Len(), ShouldEqual, 1)
						So(products.AttributeValues().Equals(value1.Union(value2)), ShouldBeTrue)
					}
				}
			})
		}), ShouldBeNil)
	})
}

func TestVariantsNoCreate(t *testing.T) {
	Convey("Testing variants no create", t, func() {
		So(models.SimulateInNewEnvironment(security.SuperUserID, func(env models.Environment) {
			sizeS := h.ProductAttributeValue().Create(env, h.ProductAttributeValue().NewData().SetName("S"))
			sizeM := h.ProductAttributeValue().Create(env, h.ProductAttributeValue().NewData().SetName("M"))
			sizeL := h.ProductAttributeValue().Create(env, h.ProductAttributeValue().NewData().SetName("L"))
			size := h.ProductAttribute().Create(env, h.ProductAttribute().NewData().
				SetName("Size").
				SetCreateVariant(false).
				SetValues(sizeS.Union(sizeM).Union(sizeL)))
			ptd := getProductTestData(env)
			Convey("Create a product with a 'nocreate' attribute with a single value", func() {
				template := h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Sofa").
					SetUom(ptd.uomUnit).
					SetUomPo(ptd.uomUnit).
					SetAttributeLines(h.ProductAttributeLine().Create(env,
						h.ProductAttributeLine().NewData().
							SetAttribute(size).
							SetValues(sizeS))))
				So(template.ProductVariants().Len(), ShouldEqual, 1)
				So(template.ProductVariant().AttributeValues().IsEmpty(), ShouldBeTrue)
			})
			Convey("Modify a product with a 'nocreate' attribute with a single value", func() {
				template := h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Sofa").
					SetUom(ptd.uomUnit).
					SetUomPo(ptd.uomUnit))
				So(template.ProductVariants().Len(), ShouldEqual, 1)
				template.SetAttributeLines(h.ProductAttributeLine().Create(env, h.ProductAttributeLine().NewData().
					SetAttribute(size).
					SetValues(sizeS)))
				So(template.ProductVariants().Len(), ShouldEqual, 1)
				So(template.ProductVariant().AttributeValues().IsEmpty(), ShouldBeTrue)
			})
			Convey("Create a product with a 'nocreate' attribute with several values", func() {
				template := h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Sofa").
					SetUom(ptd.uomUnit).
					SetUomPo(ptd.uomUnit).
					SetAttributeLines(h.ProductAttributeLine().Create(env,
						h.ProductAttributeLine().NewData().
							SetAttribute(size).
							SetValues(size.Values()))))
				So(template.ProductVariants().Len(), ShouldEqual, 1)
				So(template.ProductVariant().AttributeValues().IsEmpty(), ShouldBeTrue)
			})
			Convey("Modify a product with a 'nocreate' attribute with several values", func() {
				template := h.ProductTemplate().Create(env,
					h.ProductTemplate().NewData().
						SetName("Sofa").
						SetUom(ptd.uomUnit).
						SetUomPo(ptd.uomUnit))
				So(template.ProductVariants().Len(), ShouldEqual, 1)
				template.SetAttributeLines(h.ProductAttributeLine().Create(env,
					h.ProductAttributeLine().NewData().
						SetAttribute(size).
						SetValues(size.Values())))
				So(template.ProductVariants().Len(), ShouldEqual, 1)
				So(template.ProductVariant().AttributeValues().IsEmpty(), ShouldBeTrue)
			})
			Convey("Create a product with regular and 'nocreate' attributes", func() {
				template := h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Sofa").
					SetUom(ptd.uomUnit).
					SetUomPo(ptd.uomUnit).
					SetAttributeLines(h.ProductAttributeLine().Create(env,
						h.ProductAttributeLine().NewData().
							SetAttribute(size).
							SetValues(sizeS)).
						Union(
							h.ProductAttributeLine().Create(env, h.ProductAttributeLine().NewData().
								SetAttribute(ptd.prodAtt1).
								SetValues(ptd.prodAttr1V1.Union(ptd.prodAttr1V2))))))
				So(template.ProductVariants().Len(), ShouldEqual, 2)
				for _, variant := range template.ProductVariants().Records() {
					So(variant.AttributeValues().Len(), ShouldEqual, 1)
					So(variant.AttributeValues().Intersect(ptd.prodAttr1V1.Union(ptd.prodAttr1V2)).IsEmpty(), ShouldBeFalse)
				}
			})
			Convey("Modify a product with regular and 'nocreate' attributes", func() {
				template := h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Sofa").
					SetUom(ptd.uomUnit).
					SetUomPo(ptd.uomUnit))
				So(template.ProductVariants().Len(), ShouldEqual, 1)
				template.SetAttributeLines(h.ProductAttributeLine().Create(env,
					h.ProductAttributeLine().NewData().
						SetAttribute(size).
						SetValues(sizeS)).
					Union(
						h.ProductAttributeLine().Create(env, h.ProductAttributeLine().NewData().
							SetAttribute(ptd.prodAtt1).
							SetValues(ptd.prodAttr1V1.Union(ptd.prodAttr1V2)))))
				So(template.ProductVariants().Len(), ShouldEqual, 2)
				for _, variant := range template.ProductVariants().Records() {
					So(variant.AttributeValues().Len(), ShouldEqual, 1)
					So(variant.AttributeValues().Intersect(ptd.prodAttr1V1.Union(ptd.prodAttr1V2)).IsEmpty(), ShouldBeFalse)
				}
			})
			Convey("Create a product with regular and 'nocreate' attributes (multi)", func() {
				template := h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Sofa").
					SetUom(ptd.uomUnit).
					SetUomPo(ptd.uomUnit).
					SetAttributeLines(h.ProductAttributeLine().Create(env,
						h.ProductAttributeLine().NewData().
							SetAttribute(size).
							SetValues(size.Values())).
						Union(
							h.ProductAttributeLine().Create(env, h.ProductAttributeLine().NewData().
								SetAttribute(ptd.prodAtt1).
								SetValues(ptd.prodAttr1V1.Union(ptd.prodAttr1V2))))))
				So(template.ProductVariants().Len(), ShouldEqual, 2)
				for _, variant := range template.ProductVariants().Records() {
					So(variant.AttributeValues().Len(), ShouldEqual, 1)
					So(variant.AttributeValues().Intersect(ptd.prodAttr1V1.Union(ptd.prodAttr1V2)).IsEmpty(), ShouldBeFalse)
				}
			})
			Convey("Modify a product with regular and 'nocreate' attributes (multi)", func() {
				template := h.ProductTemplate().Create(env, h.ProductTemplate().NewData().
					SetName("Sofa").
					SetUom(ptd.uomUnit).
					SetUomPo(ptd.uomUnit))
				So(template.ProductVariants().Len(), ShouldEqual, 1)
				template.SetAttributeLines(h.ProductAttributeLine().Create(env,
					h.ProductAttributeLine().NewData().
						SetAttribute(size).
						SetValues(size.Values())).
					Union(
						h.ProductAttributeLine().Create(env, h.ProductAttributeLine().NewData().
							SetAttribute(ptd.prodAtt1).
							SetValues(ptd.prodAttr1V1.Union(ptd.prodAttr1V2)))))
				So(template.ProductVariants().Len(), ShouldEqual, 2)
				for _, variant := range template.ProductVariants().Records() {
					So(variant.AttributeValues().Len(), ShouldEqual, 1)
					So(variant.AttributeValues().Intersect(ptd.prodAttr1V1.Union(ptd.prodAttr1V2)).IsEmpty(), ShouldBeFalse)
				}
			})
		}), ShouldBeNil)
	})
}
