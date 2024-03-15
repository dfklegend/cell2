/*
 * 提供一个背包列表的抽象，可以利用背包式的结构来进行一些操作
 * 比如，各种背包显示，装备替换等
 * 
 * IBagViewItem
 *      代表一种可以显示的数据，比如具体的物品和分隔符
 *      每个viewItem可以创建对应的BaseItemListStyle组成的显示数据
 *      
 * BagViewData
 *      背包显示的数据集
 *      每次背包显示，需要组织数据
 *      由IBagViewItem组成
 * 
 * BagEnvData
 *      便于背包操作的环境数据，包括数据集BagViewData
 *      背包显示和操作需要用到的数据
 * 
*/