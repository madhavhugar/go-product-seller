package seller

import (
	"coding-challenge-go/pkg/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepository_ListTopSellersByProductCount(t *testing.T) {
	t.Run("should execute the right query and parse the results into seller payload", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id_seller", "name", "email", "phone", "uuid"}).
			AddRow("123", "seller-name", "seller@mail.com", "seller-phone", "seller-uuid")

		db, mock, err := sqlmock.New()
		if err != nil || mock == nil || db == nil {
			t.Errorf("Error while creating DB stub")
		}
		t.Cleanup(func() {
			db.Close()
		})

		mock.ExpectQuery("SELECT id_seller, name, email, phone, uuid from seller WHERE id_seller IN [(]SELECT fk_seller FROM [(]SELECT COUNT[(]id_product[)], fk_seller FROM product GROUP BY fk_seller ORDER BY 1 DESC LIMIT 10[)] AS AGG[)]").
			WillReturnRows(rows)

		sellerStore := NewStore(db)
		topSellers, err := sellerStore.ListTopSellersByProductCount()

		assert.Equal(t, topSellers, []*model.Seller{
			{
				SellerID: 123,
				UUID:     "seller-uuid",
				Name:     "seller-name",
				Email:    "seller@mail.com",
				Phone:    "seller-phone",
			},
		})
		assert.NoError(t, err)
	})
}
