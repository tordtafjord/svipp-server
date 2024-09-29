package maps

import (
	"context"
	"fmt"
	"googlemaps.github.io/maps"
)

type MapsService struct {
	client *maps.Client
}

func NewMapsService(gmClient *maps.Client) *MapsService {
	return &MapsService{client: gmClient}
}

func (s *MapsService) GetDistanceAndDuration(ctx context.Context, pickupAddress, destinationAddress string) (int32, int32, error) {
	req := &maps.DistanceMatrixRequest{
		Origins:      []string{fmt.Sprintf("%s", pickupAddress)},
		Destinations: []string{fmt.Sprintf("%s", destinationAddress)},
		Mode:         maps.TravelModeDriving,
		Units:        maps.UnitsMetric,
	}

	resp, err := s.client.DistanceMatrix(ctx, req)
	if err != nil {
		return 0, 0, err
	}

	element := resp.Rows[0].Elements[0]
	return int32(element.Distance.Meters), int32(element.Duration.Seconds()), nil
}
