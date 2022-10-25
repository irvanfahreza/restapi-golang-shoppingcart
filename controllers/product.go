package controllers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"ilmudata/project-golang/database"
	"ilmudata/project-golang/models"
)

type Product struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Image    string  `json:"image"`
}

func CreateResponseProduct(product models.Product) Product {
	return Product{ID: product.ID, Name: product.Name,
		Quantity: product.Quantity, Price: product.Price, Image: product.Image}
}

func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form

		// Get all files from "image" key:
		files := form.File["image"]
		// => []*multipart.FileHeader

		// Loop through files:
		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			// => "tutorial.pdf" 360641 "application/pdf"

			// Save the files to disk:
			product.Image = fmt.Sprintf("./project_solo/project-golang/upload/%s", file.Filename)
			if err := c.SaveFile(file, fmt.Sprintf("./project_solo/project-golang/upload/%s", file.Filename)); err != nil {
				return err
			}
		}
	}

	database.Database.Db.Create(&product)
	responseProduct := CreateResponseProduct(product)
	return c.Status(200).JSON(responseProduct)
}

func GetProducts(c *fiber.Ctx) error {
	products := []models.Product{}
	database.Database.Db.Find(&products)
	responseProducts := []Product{}
	for _, product := range products {
		responseProduct := CreateResponseProduct(product)
		responseProducts = append(responseProducts, responseProduct)
	}

	return c.Status(200).JSON(responseProducts)
}

func deleteProduct(id int, product *models.Product) (err error) {
	database.Database.Db.Delete(&product, "id=?", id)
	if product.ID == 0 {
		return errors.New("Success Delete Product")
	}
	return nil
}

func DeleteProductById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var product models.Product
	if err != nil {
		return c.Status(400).JSON("error")
	}

	if err := deleteProduct(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	return c.SendString("Success delete")
}

func findProduct(id int, product *models.Product) error {
	database.Database.Db.Find(&product, "id = ?", id)
	if product.ID == 0 {
		return errors.New("product does not exist")
	}
	return nil
}

func GetProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var product models.Product

	if err != nil {
		return c.Status(400).JSON("error")
	}

	if err := findProduct(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseProduct := CreateResponseProduct(product)

	return c.Status(200).JSON(responseProduct)
}

func UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var product models.Product

	if err != nil {
		return c.Status(400).JSON("error")
	}

	err = findProduct(id, &product)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form

		// Get all files from "image" key:
		files := form.File["image"]
		// => []*multipart.FileHeader

		// Loop through files:
		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			// => "tutorial.pdf" 360641 "application/pdf"

			// Save the files to disk:
			product.Image = fmt.Sprintf("./upload/%s", file.Filename)
			if err := c.SaveFile(file, fmt.Sprintf("./upload/%s", file.Filename)); err != nil {
				return err
			}
		}
	}

	type UpdateProduct struct {
		Name     string  `json:"name"`
		Quantity int     `json:quantity`
		Price    float64 `json:price`
		Image    string  `json:image`
	}

	var updateData UpdateProduct

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	product.Name = updateData.Name
	product.Quantity = updateData.Quantity
	product.Price = updateData.Price
	product.Image = updateData.Image

	database.Database.Db.Save(&product)

	responseProduct := CreateResponseProduct(product)

	return c.Status(200).JSON(responseProduct)
}
