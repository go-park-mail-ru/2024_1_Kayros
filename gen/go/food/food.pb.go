// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: food/food.proto

package food

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RestId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *RestId) Reset() {
	*x = RestId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_food_food_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestId) ProtoMessage() {}

func (x *RestId) ProtoReflect() protoreflect.Message {
	mi := &file_food_food_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestId.ProtoReflect.Descriptor instead.
func (*RestId) Descriptor() ([]byte, []int) {
	return file_food_food_proto_rawDescGZIP(), []int{0}
}

func (x *RestId) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type FoodId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
}

func (x *FoodId) Reset() {
	*x = FoodId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_food_food_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FoodId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FoodId) ProtoMessage() {}

func (x *FoodId) ProtoReflect() protoreflect.Message {
	mi := &file_food_food_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FoodId.ProtoReflect.Descriptor instead.
func (*FoodId) Descriptor() ([]byte, []int) {
	return file_food_food_proto_rawDescGZIP(), []int{1}
}

func (x *FoodId) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type Food struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name         string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description  string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	RestaurantId uint64 `protobuf:"varint,4,opt,name=restaurant_id,json=restaurantId,proto3" json:"restaurant_id,omitempty"`
	Category     string `protobuf:"bytes,5,opt,name=category,proto3" json:"category,omitempty"`
	Weight       uint64 `protobuf:"varint,6,opt,name=weight,proto3" json:"weight,omitempty"`
	Price        uint64 `protobuf:"varint,7,opt,name=price,proto3" json:"price,omitempty"`
	ImgUrl       string `protobuf:"bytes,8,opt,name=img_url,json=imgUrl,proto3" json:"img_url,omitempty"`
}

func (x *Food) Reset() {
	*x = Food{}
	if protoimpl.UnsafeEnabled {
		mi := &file_food_food_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Food) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Food) ProtoMessage() {}

func (x *Food) ProtoReflect() protoreflect.Message {
	mi := &file_food_food_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Food.ProtoReflect.Descriptor instead.
func (*Food) Descriptor() ([]byte, []int) {
	return file_food_food_proto_rawDescGZIP(), []int{2}
}

func (x *Food) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Food) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Food) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Food) GetRestaurantId() uint64 {
	if x != nil {
		return x.RestaurantId
	}
	return 0
}

func (x *Food) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *Food) GetWeight() uint64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *Food) GetPrice() uint64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Food) GetImgUrl() string {
	if x != nil {
		return x.ImgUrl
	}
	return ""
}

type FoodInOrder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name         string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	RestaurantId uint64 `protobuf:"varint,4,opt,name=restaurant_id,json=restaurantId,proto3" json:"restaurant_id,omitempty"`
	Category     string `protobuf:"bytes,5,opt,name=category,proto3" json:"category,omitempty"`
	Weight       uint64 `protobuf:"varint,6,opt,name=weight,proto3" json:"weight,omitempty"`
	Price        uint64 `protobuf:"varint,7,opt,name=price,proto3" json:"price,omitempty"`
	ImgUrl       string `protobuf:"bytes,8,opt,name=img_url,json=imgUrl,proto3" json:"img_url,omitempty"`
}

func (x *FoodInOrder) Reset() {
	*x = FoodInOrder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_food_food_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FoodInOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FoodInOrder) ProtoMessage() {}

func (x *FoodInOrder) ProtoReflect() protoreflect.Message {
	mi := &file_food_food_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FoodInOrder.ProtoReflect.Descriptor instead.
func (*FoodInOrder) Descriptor() ([]byte, []int) {
	return file_food_food_proto_rawDescGZIP(), []int{3}
}

func (x *FoodInOrder) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *FoodInOrder) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FoodInOrder) GetRestaurantId() uint64 {
	if x != nil {
		return x.RestaurantId
	}
	return 0
}

func (x *FoodInOrder) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *FoodInOrder) GetWeight() uint64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *FoodInOrder) GetPrice() uint64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *FoodInOrder) GetImgUrl() string {
	if x != nil {
		return x.ImgUrl
	}
	return ""
}

type Category struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   uint64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Food []*Food `protobuf:"bytes,3,rep,name=food,proto3" json:"food,omitempty"`
}

func (x *Category) Reset() {
	*x = Category{}
	if protoimpl.UnsafeEnabled {
		mi := &file_food_food_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Category) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Category) ProtoMessage() {}

func (x *Category) ProtoReflect() protoreflect.Message {
	mi := &file_food_food_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Category.ProtoReflect.Descriptor instead.
func (*Category) Descriptor() ([]byte, []int) {
	return file_food_food_proto_rawDescGZIP(), []int{4}
}

func (x *Category) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Category) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Category) GetFood() []*Food {
	if x != nil {
		return x.Food
	}
	return nil
}

type RestCategories struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Category []*Category `protobuf:"bytes,1,rep,name=category,proto3" json:"category,omitempty"`
}

func (x *RestCategories) Reset() {
	*x = RestCategories{}
	if protoimpl.UnsafeEnabled {
		mi := &file_food_food_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestCategories) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestCategories) ProtoMessage() {}

func (x *RestCategories) ProtoReflect() protoreflect.Message {
	mi := &file_food_food_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestCategories.ProtoReflect.Descriptor instead.
func (*RestCategories) Descriptor() ([]byte, []int) {
	return file_food_food_proto_rawDescGZIP(), []int{5}
}

func (x *RestCategories) GetCategory() []*Category {
	if x != nil {
		return x.Category
	}
	return nil
}

var File_food_food_proto protoreflect.FileDescriptor

var file_food_food_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x66, 0x6f, 0x6f, 0x64, 0x2f, 0x66, 0x6f, 0x6f, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x66, 0x6f, 0x6f, 0x64, 0x22, 0x18, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x74, 0x49,
	0x64, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49,
	0x64, 0x22, 0x18, 0x0a, 0x06, 0x46, 0x6f, 0x6f, 0x64, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x64, 0x22, 0xd4, 0x01, 0x0a, 0x04,
	0x46, 0x6f, 0x6f, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65,
	0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x1a, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x77,
	0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x77, 0x65, 0x69,
	0x67, 0x68, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x6d, 0x67,
	0x5f, 0x75, 0x72, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x6d, 0x67, 0x55,
	0x72, 0x6c, 0x22, 0xb9, 0x01, 0x0a, 0x0b, 0x46, 0x6f, 0x6f, 0x64, 0x49, 0x6e, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75,
	0x72, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x72,
	0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x63,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x6d, 0x67, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x6d, 0x67, 0x55, 0x72, 0x6c, 0x22, 0x4e,
	0x0a, 0x08, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e,
	0x0a, 0x04, 0x66, 0x6f, 0x6f, 0x64, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x66,
	0x6f, 0x6f, 0x64, 0x2e, 0x46, 0x6f, 0x6f, 0x64, 0x52, 0x04, 0x66, 0x6f, 0x6f, 0x64, 0x22, 0x3c,
	0x0a, 0x0e, 0x52, 0x65, 0x73, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73,
	0x12, 0x2a, 0x0a, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x66, 0x6f, 0x6f, 0x64, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f,
	0x72, 0x79, 0x52, 0x08, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x32, 0x65, 0x0a, 0x0b,
	0x46, 0x6f, 0x6f, 0x64, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x12, 0x31, 0x0a, 0x0b, 0x47,
	0x65, 0x74, 0x42, 0x79, 0x52, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x0c, 0x2e, 0x66, 0x6f, 0x6f,
	0x64, 0x2e, 0x52, 0x65, 0x73, 0x74, 0x49, 0x64, 0x1a, 0x14, 0x2e, 0x66, 0x6f, 0x6f, 0x64, 0x2e,
	0x52, 0x65, 0x73, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x69, 0x65, 0x73, 0x12, 0x23,
	0x0a, 0x07, 0x47, 0x65, 0x74, 0x42, 0x79, 0x49, 0x64, 0x12, 0x0c, 0x2e, 0x66, 0x6f, 0x6f, 0x64,
	0x2e, 0x46, 0x6f, 0x6f, 0x64, 0x49, 0x64, 0x1a, 0x0a, 0x2e, 0x66, 0x6f, 0x6f, 0x64, 0x2e, 0x46,
	0x6f, 0x6f, 0x64, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x2f, 0x66, 0x6f, 0x6f, 0x64, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_food_food_proto_rawDescOnce sync.Once
	file_food_food_proto_rawDescData = file_food_food_proto_rawDesc
)

func file_food_food_proto_rawDescGZIP() []byte {
	file_food_food_proto_rawDescOnce.Do(func() {
		file_food_food_proto_rawDescData = protoimpl.X.CompressGZIP(file_food_food_proto_rawDescData)
	})
	return file_food_food_proto_rawDescData
}

var file_food_food_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_food_food_proto_goTypes = []interface{}{
	(*RestId)(nil),         // 0: food.RestId
	(*FoodId)(nil),         // 1: food.FoodId
	(*Food)(nil),           // 2: food.Food
	(*FoodInOrder)(nil),    // 3: food.FoodInOrder
	(*Category)(nil),       // 4: food.Category
	(*RestCategories)(nil), // 5: food.RestCategories
}
var file_food_food_proto_depIdxs = []int32{
	2, // 0: food.Category.food:type_name -> food.Food
	4, // 1: food.RestCategories.category:type_name -> food.Category
	0, // 2: food.FoodManager.GetByRestId:input_type -> food.RestId
	1, // 3: food.FoodManager.GetById:input_type -> food.FoodId
	5, // 4: food.FoodManager.GetByRestId:output_type -> food.RestCategories
	2, // 5: food.FoodManager.GetById:output_type -> food.Food
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_food_food_proto_init() }
func file_food_food_proto_init() {
	if File_food_food_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_food_food_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestId); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_food_food_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FoodId); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_food_food_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Food); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_food_food_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FoodInOrder); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_food_food_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Category); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_food_food_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestCategories); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_food_food_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_food_food_proto_goTypes,
		DependencyIndexes: file_food_food_proto_depIdxs,
		MessageInfos:      file_food_food_proto_msgTypes,
	}.Build()
	File_food_food_proto = out.File
	file_food_food_proto_rawDesc = nil
	file_food_food_proto_goTypes = nil
	file_food_food_proto_depIdxs = nil
}
