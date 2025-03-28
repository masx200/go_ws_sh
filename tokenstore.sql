BEGIN TRANSACTION;

CREATE TABLE
    IF NOT EXISTS "tokenstore" (
        `identifier` text NOT NULL,
        `created_at` datetime,
        `updated_at` datetime,
        `deleted_at` datetime,
        `hash` text NOT NULL,
        `salt` text NOT NULL,
        `algorithm` text NOT NULL,
        `username` text NOT NULL,
        `description` text NOT NULL,
        PRIMARY KEY (`identifier`),
        CONSTRAINT `uni_tokenstore_identifier` UNIQUE (`identifier`)
    );

INSERT INTO
    "tokenstore" (
        "identifier",
        "created_at",
        "updated_at",
        "deleted_at",
        "hash",
        "salt",
        "algorithm",
        "username",
        "description"
    )
VALUES
    (
        '1903649975910522880',
        '2025-03-23 11:28:17.0207967+08:00',
        '2025-03-23 11:28:17.0207967+08:00',
        NULL,
        'e36421f0b7d179b82874305a4a91bd40a9cbdfb266ca987a6e5f016e19cc456d628fe127f05a9c9b11bc5af31a526887ffaa6500cae7bda92bb28f4c6e56e2c2',
        'a097f4d4c852ce0571d80eadec71d22a395fb37ce5d0ee40ce905f54c5fdf6a2ca191f1b457abe2ffdc121809687ab491bed1c376b28a58eeed4e3915cdf1e47',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1903652146667732992',
        '2025-03-23 11:36:54.5698808+08:00',
        '2025-03-23 11:36:54.5698808+08:00',
        NULL,
        '6ea66f11003367698712706cab74c71b7058d0ca3c7fa82079ff3cfe380c3a138ed24cfcdf8966f447d868b4e1cf1c9bf378f932068459046a6057ee8a50bbea',
        '7cfe723ad90b03800438ae3c976c8ce552896f1e0ec14ed2dc72e5217afae9aaaafaf19a8fef3dd882aa841d57dc78e2f795654ec5e2c916388be793779924e4',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1903659396925415424',
        '2025-03-23 12:05:43.1655677+08:00',
        '2025-03-23 12:05:43.1655677+08:00',
        NULL,
        'bdef2a0787ae3e0c7f5ae66c9e529da21ff2249f3b8dd473c0da21cdac55de18a957a96c494c3d83a3fcff351526bb5293f735a274c03e0905c2e9f40f9a493e',
        'a81da4812a66bba70de6a7e8f023ee3e7421795b05d25bfe33ba038786b4354f1ed506715cffa461c8eaa3ea481352ab5315fabc433dfe042d2dd15e56be1a97',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1903659419839369216',
        '2025-03-23 12:05:48.628678+08:00',
        '2025-03-23 12:05:48.628678+08:00',
        NULL,
        'f400a05a575140b572ffdff2ab0ecfc390c50f42240df8a4afda8b2cc3b0b5d42724d34aec8aed18d16eb9e53f32f5c1c974d5852355a1140e846f653e993252',
        'c4b3f2097b39b78f92a65d0d4381925e676983622edae2e085f8e38e665de102158301c2eb5bc50a212fa587afd6770d158af14120176cefbc85f6d844d5511d',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1903747678757519360',
        '2025-03-23 17:56:31.1949463+08:00',
        '2025-03-23 17:56:31.1949463+08:00',
        NULL,
        '82a79acda05991fac7ae2fe8c20ef72d47ab9cf43ab91251f8595e7eda8f2e3324cf5bd63b1b1fec1e1b7c908acd23ce883d94c85c9342cc565917d4c44561bc',
        '867757dacb7b48b59358b8e4b30c43e57bfa1d4ae30ef60b590696db13b046b1252ea3da37fadba7148691d77c1643eb69dc930bf7f57b8bad78dc52b165f224',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1903752786471501824',
        '2025-03-23 18:16:48.9686162+08:00',
        '2025-03-23 18:16:48.9686162+08:00',
        NULL,
        'e32e49e69cb576e056a56d0e7d966cbc4ea8456f6798c631ccf56b0a4e90557cb3674a49097b9b4d095aaac44e210bee89fc9531d92477652bccc5289ada395c',
        '81e62e2977e3b194ddedf7acdef28dee38d6d3c5536a6a54dd258cc5eae402a567085d103bbab4b25924661061f22218d2a6101824867b385b627a84bd5dd5ee',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1903755747503681536',
        '2025-03-23 18:28:34.933251+08:00',
        '2025-03-23 18:28:34.933251+08:00',
        NULL,
        '046ef72868568056cae9c6d60bc52d09ec648d328177ce8a537eee7c7437ec88f66cc5f9fc66a489bbf44428a85de74a03a492065fddbff8c93127a88df5674c',
        '77d60a4a54353034ffac3e63e80546ebeea4bde9d0c66f827dabd83680ffa805d4526f6727c22653273965192316806bdf55af1ae1b9d57b3bc220c2106467e5',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1903756731710107648',
        '2025-03-23 18:32:29.5865036+08:00',
        '2025-03-23 18:32:29.5865036+08:00',
        NULL,
        '5837bff64709d25d31f9038f89d4a854c55e35db50d43c00c3346cb9fca7b51612dfd36c93d77f7c61d161ae11c73d4d931ceb3f0f31cc2e78f2b4814d024450',
        '502b40c56ea3257efbe9c9924296e1e53ce157081f2abebe39182b9450960a3d6dd3878db3ca1ad4b7f3872e00b6216c27be6ede1cb9f78e125fb5b926aa42cc',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1903756969557397504',
        '2025-03-23 18:33:26.2932734+08:00',
        '2025-03-23 18:33:26.2932734+08:00',
        '2025-03-23 18:33:28.0990991+08:00',
        '382618873cba5741571be0eeeb35fe72330bfcc408ca0e916cbdaf4cf9eb02c0bb6cf0ffc762a8a320fd432d2c1b5c683423257df8cce9b63203039244b8e2d4',
        '4cc06756ef97bb9ec1b03a5c63954902ccad7ebffb74325835537a34c56f05e97bac077d363c78decff9c37d86a2595d019f4d39cdcdac65ea44ff720a097d6b',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1903763135346823168',
        '2025-03-23 18:57:56.3322444+08:00',
        '2025-03-23 18:57:56.3322444+08:00',
        NULL,
        'cfc90aeb864ccfe801377dbda3559e1e07af8663dce17470bb326fd1cbc038d493e923fe573f02650c59a45b59be0e4de4aaa2e410c6b17bfc519fe8a64aecb5',
        'ff0def58a618fbc7b0df437cde3987011a3e786f379b5916ed97f93eeee0ae808ddedbe478b9db8470f52a678275c326c812c225710588e2e22da052d32213dd',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1904092954482188288',
        '2025-03-24 16:48:31.3418393+08:00',
        '2025-03-24 16:48:31.3418393+08:00',
        NULL,
        'f2b78ea737ebbaea05f0396c3bfd24b9678061aaf7d24ca1af8191344abce37d6da00577b292cb41dbe2583394bdb636a568df19157aab372323d6a2f2c83165',
        'a2894933bc9fb95a29194a9830a22b433bcb6a1c803139a9d76ebacf735d750d6ae392a96b31fda94d1211d207ab31736b1fa8cc2033d72cd18be0c6231c18b2',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1904169339588194304',
        '2025-03-24 21:52:02.9705972+08:00',
        '2025-03-24 21:52:02.9705972+08:00',
        NULL,
        'f023e5216ffdddd265e328eb16a8822fb70650a92f025bfdd243d689a91c0b56abeb9662fc54be3631793f67e0bc5bea25c32d07ff3ed75ff7dc8cf8c7886a13',
        '23eca353fbc60475cf61061170237303fade269cba6278eea4bde0327cc85d599d63b0be6629394184f951e8de2632b0a3fe0286fe704d00f6e78f80413e3acb',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1904174528454590464',
        '2025-03-24 22:12:40.0926461+08:00',
        '2025-03-24 22:12:40.0926461+08:00',
        NULL,
        '6f405b6e1b9683204975cdd6cef5aa19456ddff3de83171ba20af374c411a7958ccaab8c36e36b408bc98bbcdd67be86069da9341d42b89fe90d72152a5ebca6',
        'fa7ed7bbe70efaee35dae77f02baef3349a8a20b42665d600a25f1e72f60b9ae26cdf206de2d976df46fe2bc783e03dd46ee6504fdda3e5963b8a82073acca9b',
        'SHA-512',
        'test',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1904174642515517440',
        '2025-03-24 22:13:07.2867692+08:00',
        '2025-03-24 22:13:07.2867692+08:00',
        NULL,
        'c846b19eaa502f1400b5c8cc72357ae89f7142f81221839d551598f2a6fc8a54d68df8daf0852c09b468043f15cad7a3f84c38f02c90a71e0a6cd7e1b1f354b1',
        'eaec09b838d08370c6eaf6ce27e76b6ac760890745670b7fef6534471be212135692e74ada23fd9116ef999fe741a5ceffaaf0c608ae518e4294dee480315f62',
        'SHA-512',
        'test',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1904541803678482432',
        '2025-03-25 22:32:05.3292345+08:00',
        '2025-03-25 22:32:05.3292345+08:00',
        NULL,
        'e9185ae69c02373a3f2aabd1545e892bff58b9a860bb8856b751a90c8167322728ee3f59e66b5458107cc985b75cb922f746181892f95e84ac6e36b717878a2b',
        '79deb6293ce1646d15df657f7639c7de6fe4778c37a92ccb7879e9f25d5e8209997ffaa6f3bb51132abe0299329c16ee2e89b1ba6d94b10bb17d66845c63f370',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    ),
    (
        '1904548698307932160',
        '2025-03-25 22:59:29.1376278+08:00',
        '2025-03-25 22:59:29.1376278+08:00',
        NULL,
        'a1afee26952578ec52bd9b958875e9decb571198628103d21b090619f9d5a9170e0bc2a2dd1984adaaa5516749fa64eb2a34738a4dc5431bf4b776bd170398b6',
        '922128d881796937dc5b8194dc2a8597a0c207e2c93bcf6ad19a53a003bfd74924f833481dd21be4ab43bdabb70b1dbe3dec2d254698bd9f9fc0db8f15fe2666',
        'SHA-512',
        'admin',
        'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36'
    );

CREATE INDEX `idx_tokenstore_algorithm` ON `tokenstore` (`algorithm`);

CREATE INDEX `idx_tokenstore_deleted_at` ON `tokenstore` (`deleted_at`);

CREATE INDEX `idx_tokenstore_description` ON `tokenstore` (`description`);

CREATE INDEX `idx_tokenstore_hash` ON `tokenstore` (`hash`);

CREATE INDEX `idx_tokenstore_identifier` ON `tokenstore` (`identifier`);

CREATE INDEX `idx_tokenstore_salt` ON `tokenstore` (`salt`);

CREATE INDEX `idx_tokenstore_username` ON `tokenstore` (`username`);

COMMIT;