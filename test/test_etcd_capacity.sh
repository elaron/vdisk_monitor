#!/bin/sh

for i in {1..1000}
do
	curl http://127.0.0.1:2379/v2/keys/agent/$i  -XPUT -d value="{"HostIp":"10.25.26.46","Hostname":"agent100","Id":100,"State":0,"Orig_vdisk":["2265c3e7-6476-4f87-9c9b-45ad569037af","aaed0faf-86e4-4989-9ee6-07b8e4ed3934","ccb3b9e8-9f22-46d1-9e76-7e1ba2f20d84","affe4df8-ccf2-4ebf-9aeb-c0c474ffc60c","59692c42-6032-4c4d-883f-fc9c5233dcf5","53c170b9-1621-43f4-ac1d-8e7037b4a568","0c9e858d-ac1d-4946-b340-9a82c5b9daf0","6d671ab1-3d81-4c25-9d25-9772f03692a0","1de9cfc0-94ff-4611-83c3-48aea83bf3da","3030569a-1408-4d33-87f3-56281f1f8b95","1c572b47-2f50-4438-8219-eaf7b57cc9c5","094f16a2-e044-4183-aaeb-345601bceb5d","66104958-ec2d-40bc-b6a6-01cf17f41ceb","ba7d1170-f509-496e-a10b-7f446099a8fe","15f8a30b-7531-47f7-a838-4bb41774d2ca","e6cde341-d361-4558-90a9-7ca9289fae4c","b851e28d-0f1f-4d25-89d5-bd90282a2dff","3a0a91d7-aec9-4059-82b3-d6cab38c1edd","c2f02b5f-8a61-4b79-b57e-995d710f412e","1598aa3f-f7d8-448b-8da1-433886d57745","27f670e1-cf7a-4e7a-952b-de7a03c27863","da9c8470-7ec9-4534-8020-f646298e174f","58e0579c-898c-43ec-bdad-ce6075e2a92a","73d99fae-2ab4-46e0-bb5d-671ed2c8b4c4","106ffde9-48ec-452c-896d-02bcde759b8f","0b0c9a5f-d415-412c-b6d9-dbd476650437","8bf58e51-5c1a-4b49-a06b-7ac498980ca9","ac100ac9-1194-437f-adaf-bdfd840c5261","c0f7c6b6-b088-4e3f-bf04-23fc498e77e7","dd7beda5-d000-4633-8e06-d1baf30560a8","e3a3b19b-4817-40de-86ca-112dced51cdd","ee9f1fed-7a82-4e60-95e8-4670474d68c8","ca7653a9-5b2a-40ac-b511-1d184a5228d6","45326ddc-ad97-49be-b3a0-0c0954e53e4e","f668a92f-8d57-41d2-83a3-73456b863089","71bbbff9-0052-4ae5-b02f-32931f52566a","a6278030-bff4-409d-a083-bfe80f2e3aed","a78aa31a-93bc-4213-afbb-a3464eb22e44","a6d1a9c5-a92f-453d-a0f8-b577f5cbd64f","7fd87598-a424-42e3-a442-8960b8482b46","aacb2c3a-3a91-410f-9e0d-95b27e2d19fd","0616bdfb-a183-46fd-b9a5-9065c42a265b","e472d546-9f04-4730-9b9a-f5d3e4b55aec","09ea7d29-f858-4cbf-a98b-4eb8bee9aa37","b37f2cf6-393b-4e27-86d4-8e6e9ab329ae","2db74191-5ecd-4fd9-83bf-d2381f38d3ba","35ff8f28-cb87-4e99-bd1c-d1eb143eab58","a5cd2b11-6408-4d85-81ea-0c98c16d4b63","ebcae11e-c4de-4331-ad6f-5e1dd4e0d77d","975b3031-7a83-48b7-a5f6-d7071e37ff6f","e055ac58-767b-4f8f-bd10-c993216b9cb5","a5c16ab2-bcbe-4dd2-8206-12b4da07a9e6","848d9b9a-b1a3-4f37-8c47-8531e09ac379","aee57161-d7f7-41f9-be09-8a66f3af1a47","cfb40632-e02b-4a3c-8d6c-98c631d64e0c","d5b45f45-ac4a-411b-add5-9a1be233a356","bdcbbb83-81c6-4d51-8bf5-94c02e41bc48","90b037e4-1888-4647-90ca-368b6e48d96c","b2e30342-457f-458e-b70c-ff604eca4524","167003ae-7847-4153-8ce1-8147c9c4bf20","1d0106d4-f188-4f87-b160-065c01ce6da8","f9b22f91-85e8-4e80-b20f-217d159452a7","e8492bcc-86a9-4d5b-a9a8-c9a176a0b071","fb80639b-a071-407d-a098-8317dc21c5dd","9b16660a-fe39-489e-b2bb-106fc77bced4","a23c94e7-3e1a-44f4-8d43-be54fb1f740d","6279b4c6-00df-41d6-9cad-70fee058b3a9","9d27e104-04c8-4e30-b8de-7d7f97e2f549","57a949f1-9931-4790-8e2f-919e1cf789ff","9b57a6b4-ea91-4d95-804a-f6f3cd73fa11","0230ac08-ddc5-46a3-89ce-17b049268da2","11bf162c-f66e-4e26-88ad-2033228bda62","f970299e-6c10-4619-852e-d44ed946dd28","805a6dec-da9d-4d92-8bcf-91ab7a31827b","86b52fce-587f-405a-9fa9-81a1dde3ceb3","14d988bc-dc61-4256-935d-fed54511a2ed","90c26c60-e460-4f1d-9e34-dff8abaa1f29","42dbe435-6367-4885-a524-6009fd8bebbf","93a02e00-f7d5-48df-8715-79fbb69433ee","dd7586d6-9f0c-49ee-a7d1-61633ee3fe9d","1571f598-f9a8-4f61-8c33-52b4f9008d93","ac123074-18d2-4b1e-8030-b73d4ebe5679","0ac358c6-3719-407f-bc4d-b31fa31bfa5d","d504364e-8f63-4264-a8d1-0c2727cf4a7d","cae51ead-5126-47bb-99cf-c248dae19fcb","6e0bebe9-3b1b-4f5a-97ef-e1bd6e1fdde4","66114f6c-2684-45ee-9030-7bc8ab6b71e4","fe9105cf-5263-4064-895e-3f5d8e0a4265","32198c62-f96f-4418-a6ce-dacc2d935605","9c982b4f-2ec8-4c3e-86ee-23e0c059160d","ba405ba3-00f5-43c4-bfca-6edf092c812b","059aa5b2-2635-4c3e-8f2c-e15c6cd98d0c","f1d0e79b-410c-4b23-bb03-0c26f529a32d","36c49b69-c533-4eb6-bbbb-9460a6f05829","fb24263e-46b6-4308-8d20-86ef12b8c140","30e2da97-3e63-4ac6-a623-811325faaec6","825c045a-ac11-4303-815d-5dff205dbdf2","a1bbe955-75ae-45e4-9d41-a23379222ebe","021fd060-bdcd-48c6-9cf0-c404b33ee766","fd44691d-599d-456f-bac7-a70004cd3251","7318bb5b-97f8-44fe-a9e9-594499fc6318","243b0958-d3e6-42c4-95f1-1395ff9db8c3","05e048f1-fae2-435a-b3d2-739608267f01","489d1339-82e1-42ac-9bd0-9f9e67005e79","10c16b36-5544-4ec2-92fc-5ed50250bb05","4fb5cd7f-db2f-4d1b-8472-da5721e9281e","14af4021-285f-40a6-8da2-288931af64bf","4a7bf015-a906-4360-8c46-924d788fa865","e0985491-9711-401d-86ec-a0cf62dfb602","f3912a8c-f442-4c01-80c0-7faf02d94a5f","204efd0c-89cc-44bc-85e9-910f1891d134","912f5124-b57f-4851-86df-0b442a411146","75e6491b-abe5-430a-bbfe-d2dd87d41b7e","479e31d8-d80f-4d7d-aab4-6d1b810e40d3","e8bb6150-816b-47fc-b46b-dd28f946df59","ff5b34d2-83d8-424b-b094-642b5dc5e776","700929f1-be43-4071-9ea9-2698c284adab","bab60e57-ab42-46f4-92d6-f8e2e3bad5a7","8c0d0235-c421-41b7-99bd-2b151e484fcf","6a917dc6-1fdf-4307-9a49-155de2b840a9","802ae28a-5365-48cf-acd9-fefdd3545e2d","221e1265-28ec-4652-ad30-d4d9c6559cdb","c8a7305c-1e25-4a25-b03f-823042e80342","4c7322b9-23dc-4544-bdb4-5e4ad4158c7e","338dcf86-6408-4861-bbce-1719abd61fff","6ebc2bc3-00c8-4595-abc5-65bc3732c203","fec5665e-52be-4437-8d8b-b9da9093eca3","0171b1ec-d5d6-46f3-ac6d-f1f7358569fc"],"Term_vdisk":["98360d3d-2bde-4683-9e8b-5e7a9b150d41","a5d8176f-570d-4954-8588-666c1d5f9303","c794e7ce-4ea0-48e1-9bde-37b3e7501370","da74fd23-dbba-4544-a9ff-0f84420b394b","86dc594f-b619-44bd-9edc-13a5b1532c40","00a74b68-6d16-4409-8973-f6876ef0b926","762f5254-43d4-4791-8afc-462b5dd116db","89d01a07-81c1-43f9-af5c-1515a596950a","6461d1a9-75f0-4385-8f71-d9524871123f","cac27922-cb94-4d44-a3ce-a763be2c6de6","a3bdfac5-16f8-4d72-9fb2-4151e4cf69ae","9237f848-4ae9-4b67-b566-f699145b4e28","aa9d5630-a412-418f-840b-5aa72ce53c32","62fa15db-7220-4e21-993a-00dcec23c63f","eff04680-1c06-422f-8685-5b8b3f326e35","0b57fcb5-db64-4702-a3aa-d2651b06484a","5298b9f5-f5c6-4696-994e-ba390af05bbd","bb893c49-e683-49d8-8a57-f29e6e5e70ae","3f4b5e0a-599d-4d6c-afa5-adfa59e528d4","ac7e4e47-6126-4484-a8c8-03bb9d18a567","a74ddf04-2bb8-4fd5-870a-323836edf80c","bb9cb353-2beb-4d6d-a3e9-4eb58b527d57","6b1b1651-d289-45ab-97b4-efde473476f8","2c069743-e064-416d-84d5-c8afa72f8cea","52cd3416-3b89-4667-8841-7c82670b775d","0f4af982-6829-4672-8f76-bbba769ba942","42af7ee3-fad0-46c4-a5e2-ebcbeaa66ed2","1928dbc6-3f9a-4187-b273-213050220ec5","e93d9a68-06ac-44fb-a590-4bf99e3b1386","11f24c7f-ff87-4dae-b386-7ef893081327","a73630c4-0213-426e-b727-87258f485d04","d9f86e33-9abb-4060-82d1-38bd278c201c","eb4191dd-87c8-4083-9f28-be51eb0445fe","da6de206-db67-46db-9cc2-47f7c1473623","d4ed12be-19e7-458d-bc67-f41095c30a37","064e0ad6-63ed-4cf6-8395-d8a1552cc473","d43e5bd9-b82b-4ddc-8485-e4f6928ffb65","16822e0e-8cb5-4279-8664-b8f35e992fac","6de37f72-8d46-4939-b9e8-69d6179e2246","4d396d5c-c133-4e05-a0cd-2daa8b4f70b1","cfa7ee2e-538a-460b-9638-1e93d57cdc71","114305be-bf64-454b-9391-bf37cae7d8e3","306ff67b-13a6-4417-964e-75b00e750c81","c55c1e73-2ea1-45d2-88f8-9ad9e253b2cb","d4e1a295-11ec-4464-b484-c2cc8eb88b06","e0bbd84b-d5a7-4da6-acf2-902a9636ab2c","89754140-d447-44ca-aac8-bd5539434ec3","78b29fff-8283-4347-88f6-63cb944185f6","a4fc0449-19f5-4360-8e21-a39aa41fc531","14f9785e-d768-42e6-bfaf-dc478be0a837","c81bc93c-8b82-4e21-838a-cda5dbe9d28d","21a098e5-fb87-440c-8b85-cce893398f70","a6a9e833-5a3f-446b-80e9-82be0b537c9d","82b56b5e-c8bf-4abc-8bee-4042479b6c57","8abbed4a-a19b-4777-9dfb-8e6149857682","2a0172b9-d225-4ee9-a54d-0f2fac11e861","fc9f8056-e7b1-41df-9fb2-d47a102c8c17","63d0fc89-c64f-4389-a1b7-9b8113084782","2cf6ad59-f032-4189-8464-30afb637b373","d7cebe90-cc5c-41e6-9d96-9503fbba00bf","d90e16fd-c97d-4f94-a96c-7ccd31744f17","84945b3d-e7a6-45e3-b98f-5cc11d17ea14","04a82331-70fa-419d-a13d-66777d512532","b97448b9-e9dc-4ad7-b5a2-fc41c927dcb4","bc38e1b5-d987-47c5-9719-bac7046f428e","6a22267f-7120-456e-b9ec-ca0cd8be20f4","fd75ec42-8815-4e17-a9a8-a453a1a44a80","2e7592ef-753d-4dec-8d94-a949de3dbe8f","3dc7ecca-89e9-46fc-84ab-af23ed394199","7669918c-67ae-41cd-9689-9da3e0f4eb47","1e5a8bbf-b88e-44fa-9453-b77905548c35","5b30f59c-878d-4906-a9cc-33764d1e6b93","a69e4e70-0fc8-4641-8405-398c0a966754","c58d5012-3df2-42ed-bdca-1e29c9a4d264","515a2ac7-3906-45b6-bb23-db75cef7c47c","bde06581-908d-46d4-a437-8129fe4693d8","b97a15ab-b1af-444d-a70e-a71725bb1d0a","4d975032-4a98-4aa8-9e10-391fe384a665","96cb0d04-fe2b-4b1e-b5ac-f3fcf000a6ed","3a49e3a8-ee9c-469c-a0af-154aae4bcaca","e3b8663e-6dc5-42d6-875e-f63602c7d6da","ebd73e0a-64da-40ce-8b2c-d85bbf9e97f5","66e713fb-16cc-4ff1-a7b9-cd497dd43e13","09bb7387-505a-4c1b-a31c-22301e13269f","f6bd4cb7-7585-47b0-8019-0d65a51924e0","d57af99f-2b76-49cf-9970-1e0516c3326b","800d65f2-f70a-49b4-a29b-763fcbc3cdac","8d71f438-d461-4c52-bfee-f2310300ba4f","fadc0cd9-2b4a-4779-b8c3-90637e100d5a","d9ac28b1-5644-4ed5-aee2-166333b08d9e","0df836a3-92c7-406e-bdd1-a18d4ce5fb65","b7c749b0-0dfb-4a5a-b4a6-c34fa8713c62","c204d8fe-6221-4a31-8351-55a612f43f57","c75cfa5f-27c4-431c-9c21-06b364ff3aa6","84efbebf-3fef-4802-85a1-0910d6be2a5a","3dae8374-c25d-4856-8cd4-f4cf15ff7d29","d7021ae8-7745-49a3-9467-b88877663f03","1689f8b0-5ad0-43d5-8927-a45b3f809628","fb5697fa-ac0c-46d8-91ce-18595eb3241b","36c03932-4afe-449e-a3a4-d95d97afb21b","7f523654-9243-4cc2-b9bd-44fd8408180a","990cf8eb-ceb6-4736-9983-8f11afc3cce1","3055316b-597a-4b15-9701-b70638246d04","3ce95189-f857-407d-94eb-ba49f3e207d4","a26520be-683d-4f73-9c2d-91739b750723","c51bc2e5-b71d-42b2-8d1b-45ac7c084d28","0d00db1e-a0cd-45d3-94e1-c27f6c7e1823","111f5fe6-b9e3-49be-9d5e-3e51f1406330","c1b6a547-b5cb-481c-b020-2dfb503f9ec0","74da998f-3264-4208-bb4a-d4bb93f2b36f","084bc388-0c05-47ee-977e-471dc02e2529","78ffb6b4-b04d-4cfd-8e62-c20325f9bd8b","45a8945f-9206-4ad4-b0ec-5a9a726b663d","e77aa7b3-d6de-4568-8a99-754f10f47441","7f524bc0-ae61-48f5-8d12-16faa8038751","ac07c925-db77-404b-9388-7069009b9b02","9e661743-d1c8-4c85-9885-f262369c2b4f","d4fa7816-fad0-42c5-89b7-a3cdf0e649a7","3280657d-fdef-4654-8bad-f6c94f54c826","c529df22-8dcf-43b2-bc72-5b0e10e61642","bfc818fa-56aa-45de-8c21-3890240c8b7e","2d3a8fe8-c52a-4fc8-b21e-f09c5049dde6","86086750-6c12-4a35-bb76-c5461d4bed72","7bcc5a1f-8da8-4a51-abbe-d9857b39f21f","03faca3b-8db8-439a-8d09-cee976cb2d9f","9dccd813-95c4-4c70-b219-243de9c25151","8df0e0c3-e6fc-498e-aef2-e8c64a4787d0","e81a2c79-99cc-403f-b411-1c8c25392b36"]}"
done
