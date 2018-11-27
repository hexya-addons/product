[![Build Status](https://travis-ci.com/hexya-addons/product.svg?branch=master)](https://travis-ci.com/hexya-addons/product)
[![Go Report Card](https://goreportcard.com/badge/hexya-addons/product)](https://goreportcard.com/report/hexya-addons/product)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

# Product

This is the base module for managing products and pricelists in Hexya.

Products support variants, different pricing methods, vendors information,
make to stock/order, different units of measure, packaging and properties.

### Pricelists support:

* Multiple-level of discount (by product, category, quantities)
* Compute price based on different criteria:
    * Other pricelist
    * Cost price
    * List price
    * Vendor price
* Pricelists preferences by product and/or partners.
* Print product labels with barcode.
