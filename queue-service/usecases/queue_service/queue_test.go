package queue_service

import (
	"testing"
)

type testingObject struct {
	*testing.T
}

type testCase struct {
	name string
	test func(t *testing.T)
}

func TestWarehouseInventoryService(t *testing.T) {
	t.Run("TestQueueService_GetHealthGetHealth", TestQueueService_GetHealthGetHealth)

}

func TestQueueService_GetHealthGetHealth(t *testing.T) {
	//testCases := make([]testCase, 0, 3)

	// testCases = []testCase{
	// 	{
	// 		name: "Success case",
	// 		test: func(t *testing.T) {
	// 			testObject := &testingObject{T: t}

	// 			postgresMock := mockPg.NewPostgres(testObject)
	// 			service := restService{pg: postgresMock}

	// 			postgresMock.On("PutTest", context.Background()).Return(nil)

	// 			data, err := service.GetHealth(context.Background(), HealthDtoIn{Message: "Artem"})

	// 			assert.NoError(t, err)
	// 			assert.Equal(t, data, HealthDtoOut{Message: "hello Artem"})

	// 		},
	// 	},
	// }

	// for _, tc := range testCases {
	// 	t.Run(tc.name, tc.test)
	// }
}
