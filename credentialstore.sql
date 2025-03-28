BEGIN TRANSACTION;

CREATE TABLE
    IF NOT EXISTS `credentialstore` (
        `id` integer PRIMARY KEY AUTOINCREMENT,
        `created_at` datetime,
        `updated_at` datetime,
        `deleted_at` datetime,
        `username` text NOT NULL,
        `hash` text NOT NULL,
        `salt` text NOT NULL,
        `algorithm` text NOT NULL,
        CONSTRAINT `uni_credentialstore_username` UNIQUE (`username`)
    );

INSERT INTO
    "credentialstore" (
        "id",
        "created_at",
        "updated_at",
        "deleted_at",
        "username",
        "hash",
        "salt",
        "algorithm"
    )
VALUES
    (
        1,
        '2025-03-22 15:05:36.3315846+08:00',
        '2025-03-25 23:08:53.9351401+08:00',
        NULL,
        'admin',
        'bf0f518a0039ab4b49fb855a066b2fc6aeb4d8ffde41837ebe1853179fcd65fce608cbcc7d0a18265930ace992d1627196074774bad49db77e812b182f88343f',
        '29cc37ae06d54f75b691f367e8d779bc3252e2caee857cbe694803da4173ae087e9284625ef5f8917065431d5ac7eff717d0af010ab11715b5e5cb250f74c5b7',
        'SHA-512'
    ),
    (
        2,
        '2025-03-24 21:58:09.7387406+08:00',
        '2025-03-24 22:12:51.3705178+08:00',
        '2025-03-25 22:59:19.2700647+08:00',
        'test',
        'bc8b583ab31b7b16bd54313c1c5812f77ccdaba218e0e10d079c7857fae05f773abbeb6ae199af6852a8189e6f6e2aee78b9743b1aaf9ebb9c775dfe93e1d3d5',
        'd67a15538f43b9ddc1f12ce51c18f32ed80dc957f3abb612f5cafa2335c9608b2eea346d5851fcbaf3c65e4be0707ec8fc0898566c7142c8bf07dbbe4abc7bc4',
        'SHA-512'
    ),
    (
        3,
        '2025-03-25 22:26:33.1366515+08:00',
        '2025-03-25 22:26:33.1366515+08:00',
        '2025-03-25 22:45:53.2508277+08:00',
        'abcd',
        '2cff298c70659a7340ce78631b5f8e3d3728038c4dccaabcc136d8b0ecc162277ef4ab1898c7a92149381127b42578273d87bd6a1888f8e0baa4c2e50b370516',
        'd58d5e16a0efa07b8ee5045ff03908c6271daf25df3848fc3e142b042d0e50a248e516d1218146254fd0045d71912343bfd53c409953bbdb16ffff0fb51e79b1',
        'SHA-512'
    ),
    (
        5,
        '2025-03-25 22:48:24.5477176+08:00',
        '2025-03-25 22:48:37.5520131+08:00',
        '2025-03-25 22:59:22.2222407+08:00',
        'abcd2',
        '49b9e6bec1f3a7d35438f5b1bd91b11d36190031c665985e5d517ec6ace0d83d324b4fdf393cad6eb39220f22f8c1f55d141208983a4576a11afcba136f340b4',
        '3ca01e596f29f40657019e5e42016b10bee6a19e57b36f37c6ad04ddf5ef0fa04c8952bd34d14c76470cc2bd44445550c747f03d6bfa4b4049f19f1c7253fb8d',
        'SHA-512'
    ),
    (
        6,
        '2025-03-25 23:06:52.0425022+08:00',
        '2025-03-25 23:06:52.0425022+08:00',
        '2025-03-25 23:09:59.2219776+08:00',
        'admin2',
        '7be3ecea0cefc7436caf3eb67ba2604ce5f65e06f33161931cc185e21d297333985e6994cc11512d862382bb7e66e824926f79e8bf5c8bca4766038dc3ac0388',
        '86cb3384dbad629b756d00275df3cbdcb51c9b85741a8e5415402d0d141d5146000ab7c46b4c77c2a1f9984e74db526dd95304ceeada9193dc2c2ef464147e59',
        'SHA-512'
    ),
    (
        7,
        '2025-03-25 23:08:47.7766029+08:00',
        '2025-03-25 23:09:05.1261734+08:00',
        '2025-03-25 23:09:55.9118418+08:00',
        'qwer',
        '866bfff06abdbdf86987ed0897a76003dfcda5a03c87dcd4950fe6bc99cdeaa9c17ad0480165be9b8dfb54c827bea3e28884a4c6119fe52eb4193258474e3604',
        'bfac7819ed7fcf1d4b06842d5061e231aec210b2c97e36fddf2bedff658b66780cfa79a5b95e536744b77103cf06796060f1047a4857be569bab9d650134523d',
        'SHA-512'
    );

CREATE INDEX `idx_credentialstore_algorithm` ON `credentialstore` (`algorithm`);

CREATE INDEX `idx_credentialstore_deleted_at` ON `credentialstore` (`deleted_at`);

CREATE INDEX `idx_credentialstore_hash` ON `credentialstore` (`hash`);

CREATE INDEX `idx_credentialstore_salt` ON `credentialstore` (`salt`);

CREATE INDEX `idx_credentialstore_username` ON `credentialstore` (`username`);

COMMIT;