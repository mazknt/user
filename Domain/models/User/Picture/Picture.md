## **3. Picture (プロフィール画像)**

### **型:** `string` (URL)

### **考えられるドメインルール**

- **任意**: 画像が存在しないケースを許容するかどうか。
- **フォーマット**: 有効な URL 形式かどうかをチェックする。
  - ✅ `https://example.com/avatar.png`
  - ❌ `ftp://example.com/avatar.png`, `example.com/avatar.png`
- **スキーマ制限**: `https` のみ許可する（`http` は拒否）。
- **サイズ制限**: 画像の最大サイズ（例: 1MB 以下）や解像度を制限する。
- **アクセス確認**: URL にアクセス可能であることを確認するかどうか。
- **デフォルト値**: 画像が未設定の場合にデフォルト画像を適用する。
