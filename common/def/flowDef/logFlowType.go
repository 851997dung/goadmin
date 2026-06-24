package flowDef

import "github.com/leonelquinteros/gotext"

var FlowTypeString map[int]string

const (
	LFT_NONE                           = iota
	LFT_GM                                    //GM
	LFT_SCRIPT                                //脚本
	LFT_PAY                                   //充值
	LFT_MAIL                                  //邮件
	LFT_MINE                                  //采集
	LFT_SHOP                                  //商店
	LFT_LOOT                                  //怪物掉落
	LFT_REVIVE                                //复活
	LFT_GAMEPLAY                              //玩法
	LFT_MONSTER                               //杀怪
	LFT_PK                                    //PK
	LFT_DEAD_PUNISH                           //死亡惩罚
	LFT_ITEM_USE                              //使用道具
	LFT_ITEM_DESTROY                          //摧毁物品
	LFT_ITEM_DISCARD                          //丢弃物品
	LFT_ITEM_ARRANGE                          //整理物品
	LFT_ITEM_FIX                              //修理武器
	LFT_ITEM_SPLIT                            //拆分物品
	LFT_ITEM_SWAP                             //交换物品
	LFT_ITEM_EQUIP                            //装备物品
	LFT_ITEM_UNEQUIP                          //卸下物品
	LFT_QUEST_ACCEPT                          //接受任务
	LFT_QUEST_CANCEL                          //取消任务
	LFT_QUEST_SUBMIT                          //提交任务
	LFT_ITEM_EXPIRED                          //物品到期
	LFT_CREATE_GUILD                          //创建家族//notuse
	LFT_RENAME_GUILD                          //家族改名//notuse
	LFT_LEAVE_GUILD                           //离开家族//notuse
	LFT_RENAME_PLAYER                         //玩家改名//notuse
	LFT_TELEPORT_INSTANCE                     //传送
	LFT_TELEPORT_DESTINATION                  //传送
	LFT_TELEPORT_DESTINATION_UI               //传送UI
	LFT_MELTING                               //商店出售
	LFT_SPELL                                 //技能
	LFT_SPELL_LEARN                           //技能学习
	LFT_ACTIVITY                              //活动奖励
	LFT_ACTIVITY_LIVENESS_PRIZE               //活跃度奖励//notuse
	LFT_ACTIVITY_WIPE_OUT                     //活动扫荡//notuse
	LFT_ACTIVITY_UNLOCK_BUTLER                //管家解锁//notuse
	LFT_ACTIVITY_FIND_RESOURCE                //资源找回//notuse
	LFT_ACTIVITY_BUTLER_REWARD                //活动管家//notuse
	LFT_MONEY_TREE                            //摇钱树//notuse
	LFT_PRAY_BLESSING                         //祈祷//notuse
	LFT_ITEM_COMPOSE                          //道具合成
	LFT_EQUIP_FORGE                           //装备强化
	LFT_EQUIP_REGENERATE                      //装备再生
	LFT_EQUIP_EVOLVE                          //装备进化
	LFT_EQUIP_RESTORE                         //装备还原
	LFT_EQUIP_ADDITION                        //装备追加
	LFT_EQUIP_INLAY_GEM                       //装备镶嵌宝石
	LFT_EQUIP_REMOVE_GEM                      //装备移除宝石
	LFT_INLAID_TREASURE_EQUIP                 //升级插槽道具
	LFT_SYNTHETIC_FLUORSPAR                   //荧之石抽取
	LFT_SYNTHETIC_FLUORSPAR_GEMSTONE          //荧光宝石合成
	LFT_FLUORSPAR_GEMSTONE_ENHANCEMENT        //荧光宝石强化
	LFT_SYNTHETIC_ELEMENT                     //合成元素之心
	LFT_SYNTHETIC_ERTL                        //合成元素之魂
	LFT_ERTL_UPGRADE                          //艾尔特升阶
	LFT_ERTL_LEVELUP                          //艾尔特升级
	LFT_ERTL_RESOLVE                          //艾尔特分解
	LFT_SYNTHETIC_ERTL_POWDER                 //合成属性胶囊
	LFT_AMULET_OPEN_HOLE                      //护身符开孔
	LFT_AMULET_AMULET_FORGE                   //护身符强化
	LFT_AMULET_INLAY_ERTL                     //护身符镶嵌宝石
	LFT_AMULET_REMOVE_ERTL                    //护身符镶嵌宝石
	LFT_SYNTHETIC_WING                        //翅膀合成
	LFT_WING_FORGE                            //翅膀强化
	LFT_WING_ADDITION                         //翅膀追加
	LFT_WING_ERTL                             //翅膀追加元素之魂
	LFT_WING_ERTL_LEVELUP                     //翅膀追加元素之魂升级
	LFT_SUBMIT_ITEM                           //提交道具
	LFT_ARENA_CHALLENGE                       //竞技场挑战//notuse
	LFT_ARENA_DAY_AWARDS                      //竞技场每天奖励//notuse
	LFT_ARENA_WEEK_AWARDS                     //竞技场每周奖励//notuse
	LFT_BUY_VIP_SERVICE                       //购买活动次数//notuse
	LFT_SEND_MAIL_COST                        //发送邮件花费//notuse
	LFT_BOSS_TREASURE                         //野外BOSS
	LFT_GUILD_DONATE                          //帮派捐献//notuse
	LFT_GUILD_EXCHANGE                        //帮派兑换//notuse
	LFT_GUILD_FREIGHT_HAND_IN                 //帮派货运上交//notuse
	LFT_GUILD_FREIGHT_HELP_OTHER              //帮派货运帮助//notuse
	LFT_GUILD_FREIGHT_RECEIVE_AWARD           //帮派货运奖励//notuse
	LFT_GUILD_FORGE                           //帮派锻造//notuse
	LFT_BUY_GOODS                             //购买帮派或商店货物//notuse
	LFT_HORSE_SKILLHOLE                       //骑乘技能开孔//notuse
	LFT_HORSE_SKILL_LEARN                     //骑乘技能学习//notuse
	LFT_HORSE_HARNESS_CULTIVATE               //马具培养//notuse
	LFT_HORSE_HORSE_BONE_CULTIVATE            //骑乘根骨培养//notuse
	LFT_MOUNT_UPGRADE                         //坐骑升级//notuse
	LFT_PET_EVOLUTION                         //宠物进化//notuse
	LFT_JUEXUE_SKILL_UPGRADE                  //绝学技能升级//notuse
	LFT_SIGN_IN_REWARD                        //签到奖励//notuse
	LFT_WEAL_REWARD                           //签到奖励//notuse
	LFT_ACTIVATION_CODE_USE                   //使用激活码
	LFT_ACHIEVEMENT                           //成就系统//notuse
	LFT_CONSIGNMENT_BUY                       //寄售行.买//notuse
	LFT_CONSIGNMENT_SELL                      //寄售行.卖//notuse
	LFT_AUCTION_NEW                           //拍卖行.拍卖
	LFT_AUCTION_NEW_REVERT                    //拍卖行.拍卖.错误
	LFT_AUCTION_CANCEL                        //拍卖行.下架
	LFT_AUCTION_BID                           //拍卖行.竞拍
	LFT_AUCTION_BID_REVERT                    //拍卖行.竞拍.错误
	LFT_AUCTION_BID_FAILED                    //拍卖行.竞拍.失败
	LFT_AUCTION_BUY                           //拍卖行.一口价//notuse
	LFT_AUCTION_BUY_REVERT                    //拍卖行.一口价.错误//notuse
	LFT_AUCTION_SELL                          //拍卖行.卖出
	LFT_XIAOYAOGU                             //逍遥谷//notuse
	LFT_XIAOYAOGU_EXTRA                       //逍遥谷额外奖励//notuse
	LFT_LEVEL1KTOWER                          //千层塔//notuse
	LFT_LEVEL1KTOWERGUARDIAN_CHALLENGE        //千层塔挑战//notuse
	LFT_VOLUNTEERARMYMAP                      //义军副本//notuse
	LFT_SONGJINBATTLE                         //宋金战场//notuse
	LFT_TREASURECAVE_TREE                     //藏宝洞树//notuse
	LFT_VIP_PRIVILEGE                         //Vip特权//notuse
	LFT_VIP_GIFTBAG                           //vip礼包//notuse
	LFT_VIP_CARD                              //vip卡//notuse
	LFT_FISHING                               //钓鱼//notuse
	LFT_HANDIN_FISH                           //提交鱼//notuse
	LFT_FISHING_WEIGHT_REWARD                 //钓鱼重量奖励//notuse
	LFT_FISHING_LUCKY_FISH_RANK               //幸运鱼奖励//notuse
	LFT_BUY_FASHION                           //时装购买//notuse
	LFT_UPGRADE_WUXINGJUE                     //升级五行诀//notuse
	LFT_BANDIT                                //江湖大盗//notuse
	LFT_FLOWER_REWARD                         //鲜花榜奖励//notuse
	LFT_OPEN_FAIRYLAND                        //开启秘境//notuse
	LFT_UPGRADE_ATTRTITLE                     //升级头衔//notuse
	LFT_JUYIHALL_ROLL_DICE                    //聚义厅摇色子//notuse
	LFT_JUYIHALL_GET_REWARD                   //聚义厅奖励//notuse
	LFT_COMMERCE_QUEST                        //商会任务//notuse
	LFT_FANTASY_TOWN                          //辰虫镇//notuse
	LFT_FUNNY_SPORTS_DRAW                     //趣味竞技抽奖//notuse
	LFT_FUNNY_SPORTS_EXCHANGE                 //趣味竞技道具兑换//notuse
	LFT_DRAW_WHEEL                            //大转盘//notuse
	LFT_DRAW_WHEEL_ACCUMULATE_REWARD          //大转盘爆机奖励//notuse
	LFT_ITEM_DUIHUAN                          //道具兑换//notuse
	LFT_ACTIVITY_BULTER                       //活动管家//notuse
	LFT_TRANSACTION_ITEM                      //交易物品
	LFT_NPC_SHOP_SELL                         //NPC商店出售
	LFT_NPC_SHOP_BUY_BACK                     //NPC商店购回
	LFT_WICK_REPENT                           //恶名值忏悔
	LFT_WICK_PRISON                           //恶名值自首
	LFT_RANK_LIKE                             //排行榜点赞//notuse
	LFT_JINMENG_KILL_REWARD                   //金盟战场击杀奖励//notuse
	LFT_PRESTIGE_REWARD                       //江湖威望排行奖励//notuse
	LFT_EXAM_ANSWERQUESTION                   //科举答题//notuse
	LFT_BATTLE_EXP                            //战场经验//notuse
	LFT_HOSTING                               //挂机奖励//notuse
	LFT_REBORN                                //转生
	LFT_BANK_CURRENCY                         //存入货币仓库
	LFT_COST_FINISH_QUEST                     //完成任务
	LFT_BLOODYCASTLE                          //血色城堡
	LFT_REDFORTRESS                           //赤色要塞
	LFT_PHANTOMTEMPLE                         //幻影寺院
	LFT_SOULSQUARE                            //生魂广场
	LFT_AUTO_COMBAT                           //自动战斗
	LFT_VIP_COST_DIAMOND                      //消耗钻石获得VIP经验//notuse
	LFT_FAKEDEMON_GAIN_EXP                    //恶魔广场获取的经验//notuse
	LFT_PET_CONVERTER                         //消耗宠物智力转换器获得魔力之石
	LFT_EVOLUTION_EXCHANGE                    //进化之石兑换
	LFT_EVOLUTION_STONE_CONVER                //进化之石充能
	LFT_CREATECHARACTER                       //创建角色
	LFT_REFINE_GEM                            //提炼宝石
	LFT_IMPERIALFORTRESS                      //帝国要塞
	LFT_ENTER_MGICTRAL                        //进入魔炼之地
	LFT_SIEGEWAR                              //攻城战
	LFT_PERSONAL_SHOP                         //个人商店
	LFT_PERSONAL_SHOP_REVERT                  //拍卖行.拍卖.错误
	LFT_EQUIP_DECOMPOSE                       //装备分解
	LFT_BUY_BENEFIT_CARD                      //购买VIP卡
	LFT_LUCKY_ROLL                            //幸运抽奖
	LFT_TAKEOUT_MAGICTRAL_ROAD_TOLL           //魔炼之地取出路费
	LFT_EQUIP_SYNTHETIZE                      //装备合成
	LFT_MAYA_EQUIP_SYNTHETIZE                 //玛雅装备合成
	LFT_INLAY_ESSENCE                         //镶嵌精华
	LFT_MAZE_TICKET                           //迷宫
	LFT_EAGLE_REPAIR                          //修理天鹰
	LFT_AUTO_REDUCE_PK_LEVEL                  //自动减少PK等级
	LFT_OFFLINE_ESCROW                        //离线托管
	LFT_DUEL                                  //DUEL
	LFT_SPEAK                                 //世界聊天
	LFT_CHAOTIC_BATTLE                        //无序之地
	LFT_MINING                                //采矿
	LFT_SIEGEWAR_SYN                          //攻城战合成
	LFT_SUMMON_DIVINE_BEAST                   //神兽召唤
	LFT_MASTER_SPELL                          //大师技能
	LFT_GUILD_DISTRIBUTION                    //公会资金分配
	LFT_EQUIP_NPC_DECOMPOSE                   //NPC装备分解
	LFT_REBORN_EX                             //钻石转生消耗
	LFT_ENTRANCE_TICKET_BEGIN          = 1000 //活动消耗
	LFT_OPERATING_BEGIN                = 2000 //运营活动
)

func GetFlowTypeStrings() map[int]string {
	FlowTypeString = make(map[int]string)
	FlowTypeString[LFT_GM] = gotext.Get("GM")
	FlowTypeString[LFT_SCRIPT] = gotext.Get("脚本")
	FlowTypeString[LFT_PAY] = gotext.Get("充值")
	FlowTypeString[LFT_MAIL] = gotext.Get("邮件")
	FlowTypeString[LFT_MINE] = gotext.Get("采集")
	FlowTypeString[LFT_SHOP] = gotext.Get("商店")
	FlowTypeString[LFT_LOOT] = gotext.Get("怪物掉落")
	FlowTypeString[LFT_REVIVE] = gotext.Get("复活")
	FlowTypeString[LFT_GAMEPLAY] = gotext.Get("玩法")
	FlowTypeString[LFT_MONSTER] = gotext.Get("杀怪")
	FlowTypeString[LFT_PK] = gotext.Get("PK")
	FlowTypeString[LFT_DEAD_PUNISH] = gotext.Get("死亡惩罚")
	FlowTypeString[LFT_ITEM_USE] = gotext.Get("使用道具")
	FlowTypeString[LFT_ITEM_DESTROY] = gotext.Get("摧毁物品")
	FlowTypeString[LFT_ITEM_DISCARD] = gotext.Get("丢弃物品")
	FlowTypeString[LFT_ITEM_ARRANGE] = gotext.Get("整理物品")
	FlowTypeString[LFT_ITEM_FIX] = gotext.Get("修理武器")
	FlowTypeString[LFT_ITEM_SPLIT] = gotext.Get("拆分物品")
	FlowTypeString[LFT_ITEM_SWAP] = gotext.Get("交换物品")
	FlowTypeString[LFT_ITEM_EQUIP] = gotext.Get("装备物品")
	FlowTypeString[LFT_ITEM_UNEQUIP] = gotext.Get("卸下物品")
	FlowTypeString[LFT_QUEST_ACCEPT] = gotext.Get("接受任务")
	FlowTypeString[LFT_QUEST_CANCEL] = gotext.Get("取消任务")
	FlowTypeString[LFT_QUEST_SUBMIT] = gotext.Get("提交任务")
	FlowTypeString[LFT_ITEM_EXPIRED] = gotext.Get("物品到期")
	FlowTypeString[LFT_TELEPORT_INSTANCE] = gotext.Get("传送")
	FlowTypeString[LFT_TELEPORT_DESTINATION] = gotext.Get("传送")
	FlowTypeString[LFT_TELEPORT_DESTINATION_UI] = gotext.Get("传送UI")
	FlowTypeString[LFT_MELTING] = gotext.Get("商店出售")
	FlowTypeString[LFT_SPELL] = gotext.Get("技能")
	FlowTypeString[LFT_SPELL_LEARN] = gotext.Get("技能学习")
	FlowTypeString[LFT_ACTIVITY] = gotext.Get("活动奖励")
	FlowTypeString[LFT_ITEM_COMPOSE] = gotext.Get("道具合成")
	FlowTypeString[LFT_EQUIP_FORGE] = gotext.Get("装备强化")
	FlowTypeString[LFT_EQUIP_REGENERATE] = gotext.Get("装备再生")
	FlowTypeString[LFT_EQUIP_EVOLVE] = gotext.Get("装备进化")
	FlowTypeString[LFT_EQUIP_RESTORE] = gotext.Get("装备还原")
	FlowTypeString[LFT_EQUIP_ADDITION] = gotext.Get("装备追加")
	FlowTypeString[LFT_EQUIP_INLAY_GEM] = gotext.Get("装备镶嵌宝石")
	FlowTypeString[LFT_EQUIP_REMOVE_GEM] = gotext.Get("装备移除宝石")
	FlowTypeString[LFT_INLAID_TREASURE_EQUIP] = gotext.Get("升级插槽道具")
	FlowTypeString[LFT_SYNTHETIC_FLUORSPAR] = gotext.Get("荧之石抽取")
	FlowTypeString[LFT_SYNTHETIC_FLUORSPAR_GEMSTONE] = gotext.Get("荧光宝石合成")
	FlowTypeString[LFT_FLUORSPAR_GEMSTONE_ENHANCEMENT] = gotext.Get("荧光宝石强化")
	FlowTypeString[LFT_SYNTHETIC_ELEMENT] = gotext.Get("合成元素之心")
	FlowTypeString[LFT_SYNTHETIC_ERTL] = gotext.Get("合成元素之魂")
	FlowTypeString[LFT_ERTL_UPGRADE] = gotext.Get("艾尔特升阶")
	FlowTypeString[LFT_ERTL_LEVELUP] = gotext.Get("艾尔特升级")
	FlowTypeString[LFT_ERTL_RESOLVE] = gotext.Get("艾尔特分解")
	FlowTypeString[LFT_SYNTHETIC_ERTL_POWDER] = gotext.Get("合成属性胶囊")
	FlowTypeString[LFT_AMULET_OPEN_HOLE] = gotext.Get("护身符开孔")
	FlowTypeString[LFT_AMULET_AMULET_FORGE] = gotext.Get("护身符强化")
	FlowTypeString[LFT_AMULET_INLAY_ERTL] = gotext.Get("护身符镶嵌宝石")
	FlowTypeString[LFT_AMULET_REMOVE_ERTL] = gotext.Get("护身符镶嵌宝石")
	FlowTypeString[LFT_SYNTHETIC_WING] = gotext.Get("翅膀合成")
	FlowTypeString[LFT_WING_FORGE] = gotext.Get("翅膀强化")
	FlowTypeString[LFT_WING_ADDITION] = gotext.Get("翅膀追加")
	FlowTypeString[LFT_WING_ERTL] = gotext.Get("翅膀追加元素之魂")
	FlowTypeString[LFT_WING_ERTL_LEVELUP] = gotext.Get("翅膀追加元素之魂升级")
	FlowTypeString[LFT_SUBMIT_ITEM] = gotext.Get("提交道具")
	FlowTypeString[LFT_BOSS_TREASURE] = gotext.Get("野外BOSS")
	FlowTypeString[LFT_ACTIVATION_CODE_USE] = gotext.Get("使用激活码")
	FlowTypeString[LFT_AUCTION_NEW] = gotext.Get("拍卖行.拍卖")
	FlowTypeString[LFT_AUCTION_NEW_REVERT] = gotext.Get("拍卖行.拍卖.错误")
	FlowTypeString[LFT_AUCTION_CANCEL] = gotext.Get("拍卖行.下架")
	FlowTypeString[LFT_AUCTION_BID] = gotext.Get("拍卖行.竞拍")
	FlowTypeString[LFT_AUCTION_BID_REVERT] = gotext.Get("拍卖行.竞拍.错误")
	FlowTypeString[LFT_AUCTION_BID_FAILED] = gotext.Get("拍卖行.竞拍.失败")
	FlowTypeString[LFT_AUCTION_SELL] = gotext.Get("拍卖行.卖出")
	FlowTypeString[LFT_TRANSACTION_ITEM] = gotext.Get("交易物品")
	FlowTypeString[LFT_NPC_SHOP_SELL] = gotext.Get("NPC商店出售")
	FlowTypeString[LFT_NPC_SHOP_BUY_BACK] = gotext.Get("NPC商店购回")
	FlowTypeString[LFT_WICK_REPENT] = gotext.Get("恶名值忏悔")
	FlowTypeString[LFT_WICK_PRISON] = gotext.Get("恶名值自首")
	FlowTypeString[LFT_REBORN] = gotext.Get("转生")
	FlowTypeString[LFT_BANK_CURRENCY] = gotext.Get("存入货币仓库")
	FlowTypeString[LFT_COST_FINISH_QUEST] = gotext.Get("完成任务")
	FlowTypeString[LFT_BLOODYCASTLE] = gotext.Get("血色城堡")
	FlowTypeString[LFT_REDFORTRESS] = gotext.Get("赤色要塞")
	FlowTypeString[LFT_PHANTOMTEMPLE] = gotext.Get("幻影寺院")
	FlowTypeString[LFT_SOULSQUARE] = gotext.Get("生魂广场")
	FlowTypeString[LFT_AUTO_COMBAT] = gotext.Get("自动战斗")
	FlowTypeString[LFT_PET_CONVERTER] = gotext.Get("消耗宠物智力转换器获得魔力之石")
	FlowTypeString[LFT_EVOLUTION_EXCHANGE] = gotext.Get("进化之石兑换")
	FlowTypeString[LFT_EVOLUTION_STONE_CONVER] = gotext.Get("进化之石充能")
	FlowTypeString[LFT_CREATECHARACTER] = gotext.Get("创建角色")
	FlowTypeString[LFT_REFINE_GEM] = gotext.Get("提炼宝石")
	FlowTypeString[LFT_IMPERIALFORTRESS] = gotext.Get("帝国要塞")
	FlowTypeString[LFT_ENTER_MGICTRAL] = gotext.Get("进入魔炼之地")
	FlowTypeString[LFT_SIEGEWAR] = gotext.Get("攻城战")
	FlowTypeString[LFT_PERSONAL_SHOP] = gotext.Get("个人商店")
	FlowTypeString[LFT_PERSONAL_SHOP_REVERT] = gotext.Get("拍卖行.拍卖.错误")
	FlowTypeString[LFT_EQUIP_DECOMPOSE] = gotext.Get("装备分解")
	FlowTypeString[LFT_BUY_BENEFIT_CARD] = gotext.Get("购买VIP卡")
	FlowTypeString[LFT_LUCKY_ROLL] = gotext.Get("幸运抽奖")
	FlowTypeString[LFT_TAKEOUT_MAGICTRAL_ROAD_TOLL] = gotext.Get("魔炼之地取出路费")
	FlowTypeString[LFT_EQUIP_SYNTHETIZE] = gotext.Get("装备合成")
	FlowTypeString[LFT_MAYA_EQUIP_SYNTHETIZE] = gotext.Get("玛雅装备合成")
	FlowTypeString[LFT_INLAY_ESSENCE] = gotext.Get("镶嵌精华")
	FlowTypeString[LFT_MAZE_TICKET] = gotext.Get("迷宫")
	FlowTypeString[LFT_EAGLE_REPAIR] = gotext.Get("修理天鹰")
	FlowTypeString[LFT_AUTO_REDUCE_PK_LEVEL] = gotext.Get("自动减少PK等级")
	FlowTypeString[LFT_OFFLINE_ESCROW] = gotext.Get("离线托管")
	FlowTypeString[LFT_DUEL] = gotext.Get("DUEL")
	FlowTypeString[LFT_SPEAK] = gotext.Get("世界聊天")
	FlowTypeString[LFT_CHAOTIC_BATTLE] = gotext.Get("无序之地")
	FlowTypeString[LFT_MINING] = gotext.Get("采矿")
	FlowTypeString[LFT_SIEGEWAR_SYN] = gotext.Get("攻城战合成")
	FlowTypeString[LFT_SUMMON_DIVINE_BEAST] = gotext.Get("神兽召唤")
	FlowTypeString[LFT_MASTER_SPELL] = gotext.Get("大师技能")
	FlowTypeString[LFT_GUILD_DISTRIBUTION] = gotext.Get("公会资金分配")
	FlowTypeString[LFT_EQUIP_NPC_DECOMPOSE] = gotext.Get("NPC装备分解")
	FlowTypeString[LFT_REBORN_EX] = gotext.Get("钻石转生消耗")
	FlowTypeString[LFT_ENTRANCE_TICKET_BEGIN] = gotext.Get("活动消耗")
	FlowTypeString[LFT_OPERATING_BEGIN] = gotext.Get("运营活动")
	return FlowTypeString
}
