DELETE FROM txs;
DELETE FROM iscn_latest_version;
DELETE FROM iscn_stakeholders;
DELETE FROM iscn;
DELETE FROM nft_event;
DELETE FROM nft;
DELETE FROM nft_class;
UPDATE meta SET height = 0 WHERE id = 'extractor_v1';
