Chat

Used 1 context

Dưới đây là bản dịch tiếng Việt của nội dung g:\admin\CLAUDE.md :

# CLAUDE.md
File này cung cấp hướng dẫn cho Claude Code (claude.ai/code) khi làm việc với repository này.

## Build & Run
## Testing
Black-box tests dùng GoAdmin test framework với gin.NewHandler và httpexpect . Acceptance tests dùng agouti + chromedriver cho kiểm thử dạng trình duyệt.

## Kiến trúc
Đây là web panel quản trị vận hành game (GM) cho một MMORPG ( Mserver ), xây dựng trên framework GoAdmin ( github.com/GoAdminGroup/go-admin ). Framework được thay thế bằng bản local tại ../go-admin@v1.2.23 thông qua chỉ thị replace trong go.mod — các chỉnh sửa framework nằm ở đó.

### Các tầng chính
- main.go — Điểm vào. Khởi tạo GoAdmin engine với cấu hình từ config.yml , đăng ký các page generator, khởi động framework service fusion, đăng ký custom routes, và chạy Gin tại :8081 .
- common/ — Các biến/global dùng chung: AdminEngine (GoAdmin engine), GinEngine (Gin router), GormDBList / ServerDBList (pool kết nối DB), MyCfg (custom config), Cfg_yml (GoAdmin config). Ngoài ra có định nghĩa type ở common/def/ cho các cấu trúc dữ liệu (player, log, mail, server, items, v.v.).
- fusion/ — Framework service tự viết (import là admin/fusion ):
  
  - ServiceBase — Vòng đời service với tick mỗi giây + timer wheel
  - ServerMaster.go — Xử lý signal, shutdown graceful, goroutine pools ( ants với pool 1000 và 10)
  - Tools.go — “Tổng hợp tiện ích”: helper DB ( GetBaseGormDB , InitBaseDataBase , InitServerDataBases ), HTTP utilities ( CallToDeploy , CallToCenter ), load config ( InitConfig ), helper render table, init pprof, và nhiều hàm dựng form/page dùng GoAdmin table DSL
  - timer/ — Cài đặt wheel timer cho tác vụ định kỳ
  - base/ — Wrapper bắt panic ( SafeHandler ), một số helper số học
- mgr/ — Quản lý timer trung tâm ( timerMgr.go ). Đăng ký các job nền định kỳ (đếm online, theo dõi đấu giá, archive log, backup dữ liệu). Sub-package pageMgr/ chứa logic cho các trang quản trị cụ thể (server, mail, player, log managers).
- pages/ — Tất cả trang admin panel:
  
  - tables.go — Map tiền tố URL tới generator ( Generators ). Mỗi entry ứng với một CRUD table route ( /admin/info/<prefix> ).
  - enter.go — Đăng ký các route HTML/data tuỳ biến (không phải CRUD)
  - AutoPages/ — CRUD table definition được GoAdmin generate (mỗi file cho một bảng DB). Mỗi file định nghĩa schema, filter, form fields, render cột.
  - CustomPages/ — Trang tự viết cho nghiệp vụ đặc thù (quản lý player, mail, server ops, charts, daily sign-in, reborn, battle grouping, add item, hotfix pipeline, v.v.)
- hotfix/ — Module hotfix pipeline để hot-reload data table:
  
  - pipeline.go — Điều phối pipeline đầy đủ (build → export → archive → deploy → GM reload)
  - session.go — Theo dõi session pipeline + progress
  - gm_commands.go — Gửi GM command tới game servers
  - table_names.go — Resolve tên bảng + nhận diện “bảng đặc biệt”
  - build.go / export.go / archive.go / upload.go — Các bước build và deploy
  - config.go — Cấu hình riêng cho hotfix
- adm.ini — Cấu hình CLI cho tool generate adm
- config.yml — Cấu hình framework GoAdmin (DB connections, theme, language, logging). Có chứa credentials — tránh commit.
- SQL/ — Dump schema DB để khởi tạo các bảng hệ thống
### Kết nối database
6 database MySQL được định nghĩa trong config.yml : default (hệ thống go-admin), db_global (mmorpg_global), go_manager , go_backup , db_world (mmorpg_world_newui), db_log (mmorpg_log), db_login_log (reporter). Kết nối DB theo từng game server được khởi tạo lazy/on-demand.

Quy ước đặt tên DB : tất cả DB theo server đều có hậu tố _s{serverID} thống nhất (ví dụ mmorpg_log_s1 , mmorpg_log_s2 ), bao gồm cả server ID 1.

### Routing API cho GM Command
GM commands được gửi tới các loại server khác nhau thông qua center APIs:

- MapServer → GM2MS
- GameServer → GM2GS
- GateServer → GM2GATE
- SocialServer → GM2SOCIAL
- DBPServer → GM2DBP
- AdminServer / các loại khác → GM2S (mặc định)
Hotfix data table sẽ broadcast tới tất cả API endpoints để đảm bảo mọi loại server đều nhận được cập nhật.

### Thêm trang admin mới
1. Trang CRUD (bảng đơn giản): Tạo file trong pages/AutoPages/ định nghĩa generator, rồi thêm entry vào pages/tables.go để map URL prefix → generator.
2. Trang custom (logic phức tạp): Thêm handler trong pages/CustomPages/ , đăng ký route trong pages/enter.go , và (nếu cần hiển thị dạng CRUD listing) thêm entry tương ứng trong pages/tables.go .
### i18n
Dùng github.com/leonelquinteros/gotext . File dịch nằm trong path/ ( .po / .mo ). Dịch menu nằm trong common/menuTranslate/ .

### Email
SendEMail/SendMail.go — gửi email thông báo qua SMTP tới toàn bộ user trong bảng goadmin_users , dùng credential từ MyCfg.Email .