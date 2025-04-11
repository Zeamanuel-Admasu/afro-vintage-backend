package bundle

import (
    "context"
    "errors"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
    "go.mongodb.org/mongo-driver/mongo"
)

type bundleUsecase struct {
    bundleRepo bundle.Repository
}

func NewBundleUsecase(bundleRepo bundle.Repository) bundle.Usecase {
    return &bundleUsecase{
        bundleRepo: bundleRepo,
    }
}

func (u *bundleUsecase) CreateBundle(ctx context.Context, supplierID string, b *bundle.Bundle) error {
    if b.SupplierID != supplierID {
        return errors.New("unauthorized: supplier ID mismatch")
    }
    return u.bundleRepo.CreateBundle(ctx, b)
}

func (u *bundleUsecase) ListBundles(ctx context.Context, supplierID string) ([]*bundle.Bundle, error) {
    return u.bundleRepo.ListBundles(ctx, supplierID)
}

func (u *bundleUsecase) DeleteBundle(ctx context.Context, supplierID string, bundleID string) error {
    // Fetch the bundle to verify ownership
    bundle, err := u.bundleRepo.GetBundleByID(ctx, bundleID)
    if err != nil {
        return err
    }
    if bundle == nil {
        return errors.New("bundle not found")
    }

    // Verify the bundle belongs to the supplier
    if bundle.SupplierID != supplierID {
        return errors.New("unauthorized: you can only delete your own bundles")
    }

    // Delete (deactivate) the bundle
    err = u.bundleRepo.DeleteBundle(ctx, bundleID)
    if err == mongo.ErrNoDocuments {
        return errors.New("bundle not found or already purchased")
    }
    if err != nil {
        return err
    }

    return nil
}