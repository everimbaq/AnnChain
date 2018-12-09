/*
 * This file is part of The AnnChain.
 *
 * The AnnChain is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The AnnChain is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The www.annchain.io.  If not, see <http://www.gnu.org/licenses/>.
 */


package controllers

import (
	"github.com/astaxie/beego"
	"github.com/dappledger/AnnChain/src/vision/models"
)

type InitNode struct {
	beego.Controller

	Runtime string
}

func (c *InitNode) Get() {
	c.Data["runtime"] = c.Runtime
	c.TplName = "initnode.tpl"
}

func (c *InitNode) Post() {
	method := c.Input().Get("method")
	switch method {
	case "genkey":
		c.Data["json"] = models.GenKeyInfo()
	case "init":
		models.DoInitNode(&c.Controller, c.Runtime)
	case "run":
		models.RunNode(&c.Controller)
	case "close":
		models.CloseServer(&c.Controller)
	default:
		c.Data["json"] = "no default"
	}
	c.ServeJSON()
}
