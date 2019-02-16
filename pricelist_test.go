// Copyright 2018 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package product

import (
	"testing"

	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/models/security"
	"github.com/hexya-erp/hexya/src/models/types"
	"github.com/hexya-erp/hexya/src/models/types/dates"
	"github.com/hexya-erp/pool/h"
	"github.com/hexya-erp/pool/m"
	. "github.com/smartystreets/goconvey/convey"
)

type priceListTestData struct {
	dataCard        m.ProductProductSet
	usbAdapter      m.ProductProductSet
	uomTon          m.ProductUomSet
	uomUnit         m.ProductUomSet
	uomDozen        m.ProductUomSet
	uomKgm          m.ProductUomSet
	publicPriceList m.ProductPricelistSet
	salePriceList   m.ProductPricelistSet
}

func getTestPriceListData(env models.Environment) *priceListTestData {
	pltd := &priceListTestData{
		dataCard:        h.ProductProduct().NewSet(env).GetRecord("product_product_delivery_02"),
		usbAdapter:      h.ProductProduct().NewSet(env).GetRecord("product_product_delivery_01"),
		uomTon:          h.ProductUom().NewSet(env).GetRecord("product_product_uom_ton"),
		uomUnit:         h.ProductUom().NewSet(env).GetRecord("product_product_uom_unit"),
		uomDozen:        h.ProductUom().NewSet(env).GetRecord("product_product_uom_dozen"),
		uomKgm:          h.ProductUom().NewSet(env).GetRecord("product_product_uom_kgm"),
		publicPriceList: h.ProductPricelist().NewSet(env).GetRecord("product_list0"),
	}
	pltd.salePriceList = h.ProductPricelist().Create(env, h.ProductPricelist().NewData().
		SetName("Sale pricelist").
		SetItems(
			h.ProductPricelistItem().Create(env,
				h.ProductPricelistItem().NewData().
					SetComputePrice("formula").
					SetBase("ListPrice").
					SetPriceDiscount(10).
					SetProduct(pltd.usbAdapter).
					SetAppliedOn("0_product_variant")).
				Union(
					h.ProductPricelistItem().Create(env, h.ProductPricelistItem().NewData().
						SetComputePrice("formula").
						SetBase("ListPrice").
						SetPriceSurcharge(-0.5).
						SetProduct(pltd.dataCard).
						SetAppliedOn("0_product_variant")))))
	return pltd
}

func TestPriceList(t *testing.T) {
	Convey("Testing Price lists", t, func() {
		So(models.SimulateInNewEnvironment(security.SuperUserID, func(env models.Environment) {
			Convey("Test Discount", func() {
				pltd := getTestPriceListData(env)
				publicContext := types.NewContext().WithKey("pricelist", pltd.publicPriceList.ID())
				pricelistContext := types.NewContext().WithKey("pricelist", pltd.salePriceList.ID())

				usbAdapterWithoutPriceList := pltd.usbAdapter.WithNewContext(publicContext)
				usbAdapterWithPriceList := pltd.usbAdapter.WithNewContext(pricelistContext)
				So(usbAdapterWithPriceList.Price(), ShouldEqual, 63)
				So(usbAdapterWithoutPriceList.Price(), ShouldEqual, 70)
				So(usbAdapterWithPriceList.Price(), ShouldEqual, usbAdapterWithoutPriceList.Price()*0.9)

				dataCardWithoutPriceList := pltd.dataCard.WithNewContext(publicContext)
				dataCardWithPriceList := pltd.dataCard.WithNewContext(pricelistContext)
				So(dataCardWithPriceList.Price(), ShouldEqual, 39.5)
				So(dataCardWithoutPriceList.Price(), ShouldEqual, 40)
				So(dataCardWithPriceList.Price(), ShouldEqual, dataCardWithoutPriceList.Price()-0.5)

				// Make sure that changing the unit of measure does not break the unit price (after converting)
				unitContext := types.NewContext().WithKey("pricelist", pltd.salePriceList.ID()).WithKey("uom", pltd.uomUnit.ID())
				dozenContext := types.NewContext().WithKey("pricelist", pltd.salePriceList.ID()).WithKey("uom", pltd.uomDozen.ID())
				usbAdapterUnit := pltd.usbAdapter.WithNewContext(unitContext)
				usbAdapterDozen := pltd.usbAdapter.WithNewContext(dozenContext)
				So(usbAdapterUnit.Price()*12, ShouldAlmostEqual, usbAdapterDozen.Price(), .000000001)
				dataCardUnit := pltd.dataCard.WithNewContext(unitContext)
				dataCardDozen := pltd.dataCard.WithNewContext(dozenContext)
				So(dataCardUnit.Price()*12, ShouldAlmostEqual, dataCardDozen.Price(), .000000001)
			})
			Convey("Pricelist Unit of Measure", func() {
				pltd := getTestPriceListData(env)
				tonnePrice := float64(100)
				pltd.uomTon.SetRounding(0.001)
				spam := pltd.usbAdapter.Copy(h.ProductProduct().NewData().
					SetName("1 tonne of spam").
					SetUom(pltd.uomTon).
					SetUomPo(pltd.uomTon).
					SetListPrice(tonnePrice).
					SetType("consu"))

				h.ProductPricelistItem().Create(env, h.ProductPricelistItem().NewData().
					SetPricelist(pltd.publicPriceList).
					SetSequence(10).
					SetAppliedOn("0_product_variant").
					SetComputePrice("formula").
					SetBase("ListPrice").
					SetMinQuantity(3).
					SetPriceSurcharge(-10).
					SetProduct(spam))

				testUnitPrice := func(qty float64, uom m.ProductUomSet, expectedUnitPrice float64) {
					sp := spam.WithNewContext(types.NewContext().WithKey("uom", uom.ID()))
					unitPrice := pltd.publicPriceList.WithNewContext(types.NewContext().WithKey("uom", uom.ID())).
						GetProductPrice(sp, qty, h.Partner().NewSet(env), dates.Date{}, h.ProductUom().NewSet(env))
					So(unitPrice, ShouldAlmostEqual, expectedUnitPrice, 0.000000001)
				}

				testUnitPrice(2, pltd.uomKgm, tonnePrice/1000)
				testUnitPrice(2000, pltd.uomKgm, tonnePrice/1000)
				testUnitPrice(3500, pltd.uomKgm, (tonnePrice-10)/1000)
				testUnitPrice(2, pltd.uomTon, tonnePrice)
				testUnitPrice(3, pltd.uomTon, tonnePrice-10)
			})
		}), ShouldBeNil)
	})
}
