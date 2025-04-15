package trustusecase

import (
	"context"
	"fmt"
	"math"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
)

type trustUsecase struct {
	productRepo product.Repository
	bundleRepo  bundle.Repository
	userRepo    user.Repository
}

func NewTrustUsecase(
	productRepo product.Repository,
	bundleRepo bundle.Repository,
	userRepo user.Repository,
) *trustUsecase {
	return &trustUsecase{
		productRepo: productRepo,
		bundleRepo:  bundleRepo,
		userRepo:    userRepo,
	}
}

func (uc *trustUsecase) UpdateSupplierTrustScoreOnNewRating(
	ctx context.Context,
	supplierID string,
	declaredRating float64,
	productRating float64,
) error {
	fmt.Println("🔥 TRUST UPDATE CALLED")
	fmt.Println("➡️ Supplier ID:", supplierID)
	fmt.Println("➡️ Declared Rating:", declaredRating)
	fmt.Println("➡️ Product Rating:", productRating)

	// Step 1: Fetch the supplier user
	supplier, err := uc.userRepo.GetByID(ctx, supplierID)
	if err != nil {
		fmt.Println("❌ Failed to fetch supplier:", err)
		return err
	}

	fmt.Println("✅ Supplier Found:", supplier.ID)

	// Step 2: Calculate absolute difference
	diff := math.Abs(productRating - declaredRating)

	// Step 3: Update cumulative error and count
	newTotalError := supplier.TrustTotalError + diff
	newRatedCount := supplier.TrustRatedCount + 1

	// Step 4: Calculate new trust score
	newTrust := 100 - (newTotalError / float64(newRatedCount))
	if newTrust < 0 {
		newTrust = 0
	} else if newTrust > 100 {
		newTrust = 100
	}
	if newTrust < 40 {
		fmt.Println("⚠️ Supplier trust score below threshold — blacklisting")
		supplier.IsBlacklisted = true
	} else {
		supplier.IsBlacklisted = false // Optional: recover if they improve
	}

	fmt.Println("📊 TRUST SCORE CALCULATION")
	fmt.Println("➡️ Previous Score:", supplier.TrustScore)
	fmt.Println("➡️ New Total Error:", newTotalError)
	fmt.Println("➡️ New Rated Count:", newRatedCount)
	fmt.Println("➡️ New Trust Score (calculated):", newTrust)

	// Step 5: Persist the changes
	supplier.TrustScore = int(newTrust)
	supplier.TrustRatedCount = newRatedCount
	supplier.TrustTotalError = newTotalError

	err = uc.userRepo.UpdateTrustData(ctx, supplier)
	if err != nil {
		fmt.Println("❌ Failed to update supplier trust data:", err)
	} else {
		fmt.Println("✅ Supplier trust data updated successfully")
	}

	return err
}
