package main

type mockStore struct {
	createFn  func(l Laptop) (Laptop, error)
	getAllFn  func() ([]Laptop, error)
	getByIDFn func(id int) (Laptop, error)
	updateFn  func(id int, l Laptop) (Laptop, error)
	patchFn   func(id int, p LaptopPatch) (Laptop, error)
	deleteFn  func(id int) error
}

func (m *mockStore) Create(l Laptop) (Laptop, error) {
	return m.createFn(l)
}

func (m *mockStore) GetAll() ([]Laptop, error) {
	return m.getAllFn()
}

func (m *mockStore) GetByID(id int) (Laptop, error) {
	return m.getByIDFn(id)
}

func (m *mockStore) Update(id int, l Laptop) (Laptop, error) {
	return m.updateFn(id, l)
}

func (m *mockStore) Patch(id int, p LaptopPatch) (Laptop, error) {
	return m.patchFn(id, p)
}

func (m *mockStore) Delete(id int) error {
	return m.deleteFn(id)
}
