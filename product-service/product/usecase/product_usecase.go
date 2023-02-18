package usecase

import (
	"context"
	"product-service/brand"
	"product-service/models"
	"product-service/product"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type productUsecase struct {
	productRepo    product.Repository
	brandRepo      brand.Repository
	contextTimeout time.Duration
}

func NewArticleUsecase(p product.Repository, b brand.Repository, timeout time.Duration) product.Usecase {
	return &productUsecase{
		productRepo:    p,
		brandRepo:      b,
		contextTimeout: timeout,
	}
}

func (p *productUsecase) fillBrandDetails(c context.Context, data []*models.Product) ([]*models.Product, error) {

	g, ctx := errgroup.WithContext(c)

	// Get the author's id
	mapBrands := map[int64]models.Brand{}

	for _, article := range data {
		mapBrands[article.Brand.ID] = models.Brand{}
	}
	// Using goroutine to fetch the author's detail
	chanBrand := make(chan *models.Brand)
	for brandID := range mapBrands {
		brandID := brandID
		g.Go(func() error {
			res, err := p.brandRepo.GetByID(ctx, brandID)
			if err != nil {
				return err
			}
			chanBrand <- res
			return nil
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return
		}
		close(chanBrand)
	}()

	for brand := range chanBrand {
		if brand != nil {
			mapBrands[brand.ID] = *brand
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// merge the author's data
	for index, item := range data {
		if a, ok := mapBrands[item.Brand.ID]; ok {
			data[index].Brand = a
		}
	}
	return data, nil
}

func (p *productUsecase) Fetch(c context.Context, cursor string, num int64) ([]*models.Product, string, error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	listProduct, nextCursor, err := p.productRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	listProduct, err = p.fillBrandDetails(ctx, listProduct)
	if err != nil {
		return nil, "", err
	}

	return listProduct, nextCursor, nil
}

func (p *productUsecase) GetByID(c context.Context, id int64) (*models.Product, error) {

	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	res, err := p.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	resBrand, err := p.brandRepo.GetByID(ctx, res.Brand.ID)
	if err != nil {
		return nil, err
	}
	res.Brand = *resBrand
	return res, nil
}

func (p *productUsecase) Update(c context.Context, ar *models.Product) error {

	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return p.productRepo.Update(ctx, ar)
}

func (p *productUsecase) GetByName(c context.Context, title string) (*models.Product, error) {

	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()
	res, err := p.productRepo.GetByName(ctx, title)
	if err != nil {
		return nil, err
	}

	resBrand, err := p.brandRepo.GetByID(ctx, res.Brand.ID)
	if err != nil {
		return nil, err
	}
	res.Brand = *resBrand

	return res, nil
}

func (p *productUsecase) Store(c context.Context, m *models.Product) error {

	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()
	existedProduct, _ := p.GetByName(ctx, m.Name)
	if existedProduct != nil {
		return models.ErrConflict
	}

	err := p.productRepo.Store(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (p *productUsecase) Delete(c context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()
	existedProduct, err := p.productRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existedProduct == nil {
		return models.ErrNotFound
	}
	return p.productRepo.Delete(ctx, id)
}
