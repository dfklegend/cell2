namespace Phoenix.API
{
    public static class CollectionBuilder
    {
        // 创建所有分类的API集合
        // categorie:  category1,category2
        // serializer: 参数序列化器
        public static APICollection BuildFromService(string categories, Serializer.ISerializer serializer)
        {
            string[] subs = categories.Split(',');
            APICollection collection = new APICollection(serializer);
            var types = APIUtils.GetAllClass<IAPIService>();
            foreach(var type in types)
            {
                if (!isMatchCategory(subs, APIService.GetCategory(type)))
                    continue;

                var name = APIService.GetName(type);
                collection.GetEntry(name).AnalysisType(type);
            }
            return collection;
        }

        private static bool isMatchCategory(string[] categorys, string category)
        {
            if (category == "__system__")
                return true;
            foreach(var one in categorys)
            {
                if (category == one)
                    return true;
            }
            return false;
        }
    }
}

