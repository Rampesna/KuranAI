#!/usr/bin/env python3
"""
Test all 25 questions against the test API.
Usage: python scripts/test_all_questions.py [--base-url http://localhost:4551]
"""
import json
import sys
import urllib.request
import urllib.error

QUESTIONS = [
    "Bugünkü doluluk oranımız geçen yılın aynı gününe göre nasıl?",
    "Önümüzdeki 30 gün için doluluk eğrimiz hangi günlerde kırılıyor?",
    "Hafta içi ve hafta sonu doluluk farkı ne kadar?",
    "Doluluk artarken ADR (ortalama oda fiyatı) düşüyor mu?",
    "Doluluk düşükken talep hangi segmentten gelmeye devam ediyor?",
    "ADR son 90 günde nasıl bir trend izledi?",
    "RevPAR artışı fiyat artışından mı doluluk artışından mı geliyor?",
    "En yüksek ADR hangi tarihlerde ve hangi kanaldan gerçekleşmiş?",
    "Fiyat düşürdüğümüz günlerde gerçekten daha fazla oda satmış mıyız?",
    "Rakip ortalamasına göre bugün fiyatımız nerede konumlanıyor?",
    "OTA, direkt satış ve acente oranları son 3 ayda nasıl değişti?",
    "En kârlı satış kanalı hangisi?",
    "Komisyon oranı yüksek olup düşük katkı sağlayan kanallar hangileri?",
    "Hangi kanalda fiyat artışı doluluğu etkilemeden mümkün görünüyor?",
    "Direkt satışların toplam ciroya katkısı artıyor mu azalıyor mu?",
    "Ortalama konaklama süresi (LOS) son dönemde uzadı mı kısaldı mı?",
    "Kısa konaklamalar mı uzun konaklamalar mı daha kârlı?",
    "Erken rezervasyon ile son dakika satışlarının kârlılık farkı nedir?",
    "İptal oranları hangi tarih aralıklarında yükseliyor?",
    "No-show oranı gelir tahminimizi ne kadar bozuyor?",
    "Önümüzdeki ay için fiyatları nerede agresif, nerede korumacı tutmalıyız?",
    "Talep varken fiyatı erken mi düşürüyoruz?",
    "Doluluk hedefi mi, kârlılık hedefi mi şu an daha mantıklı?",
    "Hangi segment (kurumsal, leisure, grup) bizi aşağı çekiyor?",
    "Eğer bugün %5 daha az doluluk ama %10 daha yüksek ADR hedeflesek toplam gelir ne olurdu?",
]


def main():
    base_url = "http://localhost:4551"
    if len(sys.argv) > 1 and sys.argv[1].startswith("--base-url="):
        base_url = sys.argv[1].split("=", 1)[1].rstrip("/")
    elif len(sys.argv) > 1 and sys.argv[1] == "--base-url" and len(sys.argv) > 2:
        base_url = sys.argv[2].rstrip("/")

    # Login
    try:
        req = urllib.request.Request(
            f"{base_url}/api/manager/auth/login",
            data=json.dumps({"credential": "manager@intelor.com", "password": "admin123"}).encode(),
            headers={"Content-Type": "application/json"},
            method="POST",
        )
        with urllib.request.urlopen(req, timeout=30) as r:
            login = json.loads(r.read().decode())
        token = (login.get("data") or {}).get("access_token") or ""
        if not token:
            print("Login failed: no token")
            sys.exit(1)
    except urllib.error.URLError as e:
        print(f"Login failed: {e}")
        sys.exit(1)

    results = []
    for i, q in enumerate(QUESTIONS, 1):
        try:
            req = urllib.request.Request(
                f"{base_url}/api/manager/test/openai",
                data=json.dumps({"message": q}).encode(),
                headers={"Content-Type": "application/json", "Authorization": f"Bearer {token}"},
                method="POST",
            )
            with urllib.request.urlopen(req, timeout=120) as r:
                data = json.loads(r.read().decode())
        except urllib.error.URLError as e:
            results.append({"q": i, "question": q[:50], "ok": False, "error": str(e), "sql_ok": False})
            continue
        except Exception as e:
            results.append({"q": i, "question": q[:50], "ok": False, "error": str(e), "sql_ok": False})
            continue

        ok = data.get("isSuccess") is True
        sql_results = (data.get("data") or {}).get("sql_results") or []
        all_sql_ok = True
        any_error = None
        for sr in sql_results:
            for ex in sr.get("sql_execution_results") or []:
                if not ex.get("success"):
                    all_sql_ok = False
                    any_error = ex.get("error") or any_error
        results.append({
            "q": i,
            "question": q[:60],
            "ok": ok,
            "sql_ok": all_sql_ok,
            "error": any_error,
            "request_count": len(sql_results),
        })

    # Summary
    ok_count = sum(1 for r in results if r.get("ok"))
    sql_ok_count = sum(1 for r in results if r.get("sql_ok"))
    print(f"API success: {ok_count}/{len(QUESTIONS)}")
    print(f"SQL run success: {sql_ok_count}/{len(QUESTIONS)}")
    print()
    failed = [r for r in results if not r.get("sql_ok")]
    if failed:
        print("Failed questions:")
        for r in failed:
            print(f"  {r['q']}. {r['question']}...")
            print(f"      error: {(r.get('error') or '')[:100]}")
    else:
        print("All questions: SQL generated and executed successfully.")
    return 0 if sql_ok_count == len(QUESTIONS) else 1


if __name__ == "__main__":
    sys.exit(main())
