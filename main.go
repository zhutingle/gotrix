package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/zhutingle/gotrix/checker"
	"github.com/zhutingle/gotrix/global"
	"github.com/zhutingle/gotrix/handler"
	"github.com/zhutingle/gotrix/weichat"
)

func test() {
	fmt.Println(os.Args)
	fmt.Println(filepath.Separator)
}

func main() {

	bs, err := ioutil.ReadFile("G:\\workspaces\\gotrix\\src\\github.com\\zhutingle\\gotrix\\func\\gotrix.xml")
	if err != nil {
		fmt.Println(err)
	}
	bs = bytes.Replace(bs, []byte("GOTRIX_ENCRYPTED:"), []byte(nil), -1)

	decryptBs, err := global.AesDecrypt(bs, []byte(global.Config.Args.Password), 256)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(decryptBs))

	//	decryptBytes, err := global.AesDecrypt([]byte("F0:^0^gg0@32A>CD\`GPA9F0N@HNHN_[=dmOd>7[L6^=SWgY3EB`<145aI=Qm;:c9RK\[W5=_nXXgeUS0F;R677giKWd5gmZHT@4LJgZaZ:oVl7RQ=b[^Mg[CQ<bm`27b=NHjPgFUP7Ph?COSgh:mMD=03MLAK81U1m0^@@hN7@i=;dB:7ICHSc4`8Y2\=>bA`cE`V`o:mNc=F?AmfOH^IHj:f24fbkNKMNZc4a?7L5o]a<5`\;FSihG^>D31K]J:j\TdoVf:RO7VU2_H1fKkNWjk_[57P`6TXNVCa`Ao:Zn69<kPYHR`gE9RFOMYf57e@7MWY=HGVoAVe>C]hcdHO[lbiQa5S81MAa:Khf^bHBaY57MjK=c^F:8XFFY324k\<l`jFCI5kfc9=C9k18ag\>@a=Z[8VN02Qdhg_dh<nR_iNkMAW\k5Oo]F^i@YXjR;HP<?b^N[FRhfTZ<d]Q?GPZL>f[J>1RZoAK^14lbTD?BgiW8amA5:;<`GC1ZVXSKIBWN>i?CU6fNJaNFB]kD9i:8lRVPaRF=2VACc13TW3>gFBaGiSCA2D8;gAQl5BGfXj7<U@PYX7V9X0=8D62kRKP7hghDWKAS96aEOXOTakbcT;GmH?hm\F___iTQC2AeCOId\4LVGG7n@?26bePV@R@V==8^IIVee1CX6_m7bZ[WZ_;B@fnYCnffBn5?k>LJjIQ7`DWNZihW9fmG2R84@Kl;fA_eff?E03?OVhL5T6hCWK1:`f3P=U000=cJ[BHA2Gf9MlFk0^WKH[TF7CKaNmT\e[`>^iVlG]9?AbVG0_RefRol7@;QbC3=I<:;eg>FI15GJ4>08D3R=2i58kmR`_O0Y]2]E\PYV`=Hh5PZgSX7cgRa5?bbEgk<YGjljL2cdU_KDR@je7DgbmAC=GE2R_BZc<iHUSjnS15@fQ5kO@40iSo8oK2eFiFVJD;9eAKoD^61liOfA4VK7i[5ZiaHZ<1c4K]P54_X4=YJ6S]APoVHB2gc\XF9m7k7Mi2C<d^;@9AVWFY>N?PL^FJc<mGOn=;o7V7RR0m>;7gJ_M>S48Vg4\2;K0Vf0E2[G8Pb0E7AcFZ1kW9hXX_96:EH1l<`4Bch<?H]4Xn`HFc5O8j^e>G59[`H6a7@8^?1B9V\62gH4NoefiV=\UK4L\9m:]nGJb_WE?NDK2JdkF9U1JLIem[Yob?3D4j7]?79ISgHSjjVYd4^58\40eDRS?Z4`XQ[Z6RO5j[DlE08d@Wb;:l@0STle>La946li[;9e6MRGe\lWZfOB[lP]WXMJ^:0_Sk\V_FBZQKVA68TK9jn_[P_nV1oLVAPC;Bib@4F2ke9ZO:nWMFckZ^j^^5ICLG:k8\??76^jbPD`4l:U;<aOP=<?b^m=9e_C72jL^@k\^kb`H_H<F3V]jAoS>VJ@h^Sg_mC9Xg3QCdBLGdTTD_[A@YFifOKh1NLFa4cT9IRAF0;N94l[0o<2l2G=O@Nijec5QE>KG7JK4dOiCn1\Db2T]GZ?7Z1iNgKLOKSGYTXAXF2V]SZcAR7eCT^AejeUHJZ03l^SJ?_a=Ba_ba>kb9hmBEG<G@K_O071MNX7YQocb;QE@Q53DgjafXGj]I3SB[]\P2:ka^AUdD27eA\3Jg\I:FF4gl`>kUK1c1FnA8hIdIf]>E=4U;CLM4WcSEaahH5d0imV[B0C5lTBlYUXWcKEbKn4iZd:kS:YTb<e\RmCRLe@MfBIRj4E@eiRh9@dMelRQcD_Kcb8=C?0V2kQJj\Q;eBHi[9DcF3E0E6dkl>5]QXXZmWkb\emDk=GYOWLdakKeNJL1d3JJLNnhmAV@k0CDS[5FfX_Z72]<VVGMQZ5cF=6lnVS2<nanc@QBm_HW_IQeMMFjF4kM4RG<`^:J3@5ZI[HXS>6Whj<hcN13flODA_6?X^J1W=He[>5Z70egE7YJbP<c7R8^;>]C@lOodkfVAFGCGm57B2e18T7_FR>Q@_C^dl>9[HXQ:_=eYoa=8cZhKL><^`9b@h>5EbC68b@H4U`RI:LkIk7d`>ZS7<=[\5@4MN;o;gJ6ClZ:58=jWC]<PkaG[dkX^L<nGS=Z2HNNJ]lOiSi=EITZa6Wm?0O;?gPaoe=FJPZ2m9Z]X>omI];k:j@j6LIEb=Db_doE2XZ>ICYd>K;jb_VPBP[gNSckWch;51Ojl9hR7]\S6[E^PPdLe4Qe3M>eee]IKV9PJAN=dNI]_\eH@>_FKnKa3UW1c^Q>9de9EW\1I4o0_k1Nk<AVZP\=Hj5j\o:OUHGQIGhO[bb:abCKWRL8AXaC:oTVjU85ljNYJbT1>Wbe3PgNO>3AlgOiH`mkdeOb[7a4Kk0MS_aU`NV4bG\fC?O`>VTdUh6NCV3UU8kjH0YOnTmok_=1i<UolLbo3SAA6gBiBXE?a4[7\KKD8T0QGRdbE7FH?Z[3K`kTB?QOQNhdoPccH:HCm5Soa?_[Va_h\gj_52UH\mQB[b<0HkPT]7aNZ=[EEaJI\`i^=;mLP?F]Q1LL7Q;oAeG`f`BceNZKK;A]DNRI7SXaI==1oPVR4O:QHZQnc_bJBB1?C^EcLkKkAnoPV2PkF=kZm30?Lbe>j6Ye[cT`25RZALAd0SS@ZU^9VC3lgY?5[b:^gU_85d3KaWC>661:bAY]_oK>Z>V1=^ZQDn8FGM@`FF]Di1>9M6Y`9\4DNZd4?7UG9RCOeURFgJila9Id8X9>0CP^@>3XTFeClcDmL@eFK]diQGnJi2bnG>5facg3<oKHbO6^Q7U2J@>BM7D]IMbJeO;aJ]^6_CWbL\mTBNBbalKfc_gf=6_`DTB@XWeLM7i\@7ThaZUkoBR7PJn`J93H6>79_>EiQ=6@UUCX[ON06JcJk5QOHo4PLAO6F]D8Uc1FaIWk8j2Qi_o<1n_TnaICcNmES?P09LY@_H1oSTjXYCCSG6VCkJ0QeMXhTMBF]YjV>T^@ho[CWFo\I[m1[:<IM1]AlV0?oWh>g;]fIlIkkIjl9hIRfgJXUWgGb2]^hkk9CUgLG;cBmDfn4ggdodn\dd6EaR8kHAkdUF@YZlZi[dLc<fc=jgFl69PkP0_K2P_Ck3J_`ZNg1:Udghh_oS1dh`bbBf`c5R9O@gbT<ooOZBS2TG_WJ\<anS_[S>=^nAM7;BKd?^7iC^=5k42o_];eJfT_oR?ZOked4gJ4SWiRo7d>^9hQ4[]_m]]FE4jCEgZn4Le@F1I^[FhR<?5DEm??Zkg]g]1OOLOdCQA_A9PR_i3`DF<5dlS_[ZZViEES9mlED;7;8[1]>oLA3oBdJBfG7C\4?2fo=og>`N8QcG_AR@6FQfSRVH^a@Je[Y>_n51l\L]me=:N[oNkJGkjeJ;GP5eIR@?eLmPeS\]Kb[S22kAP<eAHVdo;QIRTDN=X;SLiafN7_3aI0k[`n8Wjgd:05<;NZmHMV21gb4j@YLaZn6]5L1:I@3jFB_Yk?TEXKkC6f2QjED8kk5@TKkFiOcg:WoW7Ge[e6Uo@SZ4Mh7TbJ6^R<J1]4C8Yb8b;\QN?Eo1hWAjj7i20j5nHXcH5G]VMWk8QLNIn[omnJ?7hn7WQQnA?Oa<^bGC5M1TD2Uf\f0=g@::I>e@oJ;ZZjJm80NoGbigcVhBm0Eb`8`T173B1h8[IC\A9T4Hld_7Qb3QW;PX=B3EIl<G:9_FoEPK[a<RXk:Q^5:cna=LS^O:MW<dggQ16M6JFQdL?=eDgW<30d;Cb<ejJc[m]_PNGDVfGJ`H4ODc6?RAoEFJ]<YUZ@\IRnkX?nVJ5PbM1<EAoI7W4T7LQ1WR?Zkg8kIkdc`Eemi[[1HkZM3o;Xmb_CN@i;_m;4;]oTLNmgH4l>Me\\95E@leD]jX]5D[TEEh?RajE?l_;\ETLOGHIZm75n=>[DUVmj1n[BVC[mfUdN_oQJU5T9FGKRNUn4cZo<d36kU]a]cXOE=1DgZlA3;NgIRgl>;;7idJ7PQ2;G=]:YfaS87DX5ZAVVlRMJ2keIDHUbADbcQCD2N6R\<;em6LV5eS>HIbK@WHJSmhZT[AmJfkAg[K]BOY;U][`2SU<]=bH3cG9PRF8m:i2jE8_bX5c8V=\]J^=nL2`7nFa7[nMUM>^FV`0gIe=6[L2Tk:;Lh`]Ff;Shn1\l8WP:X6bRCdKA:647jIDajVSjKT6HTiO3cO<=I9^BNWo>Skbl0@l8:E7JQ1:hccn3;JUX^1aGkC7FO2N8HC9Fc[FKE3RFEVWTW`4J24OG9VKAmbl8QZh8088^VeG[b8EgEei@nD3?0[;XmNXo[Zii^\W2bLkGeVM6n?oA6kRQ]YYDW=_YUWhSOgI;\LV;05UILKmbWO2=E^P@ZP>mEen:7[n2J;1]CK8SKYi>V2U^]Lf:45AO;V[3`g45miN<A[L30kbOF[IQ;\OgXceKYXJ0[NVQ`QONi@GNVB\NLLkVTkFHMDR];oj9`eTeR`lIk1fIU^cU?Cab?6g`>MP5OfUe70Cao]Ehc;^4eDdZcO\f_<H<7Y6S2B0IcfZ6mDn@WPflVI8kNhT3m6VimHH<>O[ob2eYcDWU>o>SJMEUDMSV1=V3gFVIYIX;[5l3giSSY5H\ZJ_K>YbI2L]aF74[JB_5BSkb=B_HV<DdW^[RbYoUJUbF0TCh^J1GaW7WNOAONHWRmaO2EfVlV4dFY0RG=e_4e>oDim^iJge7?KF6Z5Ab]O`0j>Vj8M@R:J1c>`IV<\SnXA?o\C2iEWOZWK@DOm6P3ollSEcMVlm`fmQfkWhnoeFa97a:^0JFBSRRYIQBB9b7]SgJAienKJNME`GGInnMN`:9OPkB41><1jHKAfUT4]1Obn9;X\N`D6??[:CFAg[IZ@NdLTFAHBglR9d8ocS?NMm\gIX@T6mgFYiKXbAF]hZ7m5cCh@6g\R7TUP[0H1Zo<Z_U<iAeNBfL<mN5_4H=mA<8oJOKanDl_Bfi5hCR5DAKSlidemNFN8LA;\7c4Aji;d36:N3`ARBXo:hPZX9E^AXKM]1U@GW83gdPfg1aekUNfScGYW`?G2_[:6TKPgRQ9C8m@ObF]\MK6N5CRM<`B;lSZn<0km:0g^`GmjKbnd_2I9DlLaXD;06A2OUd;0KbWgZ`0g9D:o\ClQA2miG4PA\m?P0BTmhjRkHESb4lc0ECW84oRo3jok[oAaJIM5K4415bmcLkDJ8_8VhkGEg>_?dOl6^jo`nHVaI@IG0LDC^f<ZPdBE;dU\L4XAni[hmM`go>03X`6DC^dI2cW4MVRbhNHL2>2>;>YK\1egFHfUeVINH[Y@0K84B1:cng`CG6`?_1f`KZmL:=2blCe0EXfDkf>6[5:ZR9jG6lj4E;5m1;CK^99Rcc_e5SE12<QWkcVaFnZ]WeWj0CJc3IH=SdE8OfXL:@X7P\biin:9DY:bWhUajhZ=TMR[DS@bQ9=7IHZ1aOcKbdb75hYC8R3L>UNo>nkS\eGCKAbKGMmjEF8d?;9aZISGlU6gnKDQFK^<@;@WGBXFj[BSm?o[09RdI1bmGUKfgN?S?3TOo5KalK1FAYi@205ib?W>2TZ[h8n]L?ILO4]5e[M`39KCcg3k^@M2nX`kc2ga>bgO@6?>3;M7ZceI^FfTBFBVX[6H^oKTIe;24id=6WjgOaciMS1eE1Zf_k<F`]@Qh8>d1ICV68REkNPA0Z:b_ZCUImaWm;fM^EWIbD=22Oe50`gBJ<AZ]nUa9EGW32\mM[?5=XfgE[DKkMk@10_c5i?\[K8Oho2hnh65U;FUf^j_j9=A[RYZX5]lH^Bn65Xj3J9>9jBoWjCY;;B_USmPKVILK^1V9dP0mRT^Fioo<bCdLB;PjC1aRPnP4M<DY3S5a=UmhOHh=fgOEmeGDO<2SIHJ\cQl4@>E3GAR`c?Wb9HL[0ndlO[G1]>LOi:E`[?CYl8<I[[EA1X[cO;Jim;6:h4k;gYj?<gN?]^ZRhMG`NMB?>NFPL;6abmQ9;AOF<QD4HkX03\_Zca6g=S8_lhRiSomfadDa3HEfFdTT^9bZgASU<l^EP0^\D;]Q5>6<@>2Y0hNRO0TZ4QeH[@Sb]9Q\;<cJ6ZE_Q>c3m9Ddimh20[Y^I9a[;Q>@Q:VH\:fhG=P_iZ<:73PdLTjT]?AHc63_Un`dG>UF3SGAm^oIWD5YCFV12R4SnD=GEb@>L9`eZWl`3?DK]Wn571@PHf4K<YYecc___hYc7<=`ObNhBa2g`Nj]TZTd@]F0f\2ATbbo6mON97c`5LnIjaDkdRc6j\8ejQCcfG;BE\@7AMT@@9ZG`=DPG;R84MIRhVkgBZC6XjRBT=ogfIo6jjfPaUKj2<WIM7[UNV;?nDn[0<U0]m>cmL]QTcOk<k0OodFjajXDdJ`Q^CSHK6lmO=i[F3Lm`JH;]dK;[GXhfO4`nE]5SW8_L:md1S2Q:lkVM4R4O7Aj7K;54=b^=8Ca6lkgRa^POBLBi?2?Go_54;=>eNSZmFh2[2<3@IKe]EFQ1l]nRgQm:fIW][ZMNK?d:743DHBn\H5>DGlKZ2Xo?:5lc@PhAf^hXI`6YNP92ZW15KM0i;cTU6nMAWOhlJZdOITY:ne72`^3[EKg6Ba^0=cjg4AAUUbIokFEbMXQ]4XWNm@d1Z_iG@H]3F1[TZEZN4FgjQYZ9o[NFLAMVb6fAK6ATYKI2cldB@oh6o020=UHA@5=Zff^K\G:ig:EBDO^UboGFNm1EXGK\3cdkO_[B5Ra_IVZi3?286@2TJi<^O>Aa`:ON>JIohSle?F`OS[dhYAJD;4o72i2jOX`O8f?4cVU=1FW>=UXmQle9?5f1??<kgf;BAmGNJQ8cHHMCRaU<ZS;BbcSFR;`B4eOLe7[]`Fjo81>]7lM00<?f:@]GLSACVcW;hASl8kOgo40[4^A67E_6jl35WX@:a<WL2bU?iLWjI_4HTFo=fen@8[P9>@hgeA`S4Oh50`:nJ3dDEF`^@IfaEX@3[4`9P[Lmg\MYEB5DEXf<6L_SOe5Y@2\kaJ@Yc2@AE`=HbDZ[TUD\O8:1Vl0dSZmg\3WN7SWYKR3Jh3eMoK;6A<eQ=WTGbH2YT\gCR<^B7MO<i>H7G`RoJF]^5LmTWDl[HW]U2PiDKOX_\USFLi?AWkll3kN1OTe77_C;TI=9FfVilJ9]a<XZVS>2VNFECV>2Xj\YmITiaDZ[ljDAX]hLG`2@4GV:]2MLCJ=V=^:NgGLP[<dI4oH0aK92VHeWDJ[3hEk15Zm1`UI;PLfChD:cN]2l_neEG=f2EB5BWLRFP?87A_Ro`1H;\WEjgNZ^@><DJBIkK1`\5SidMhE]In3B0gNHJ6an]<HHeX=[CjWXN>@D9CM9U;?gcV7IGgA>kJN6BfN^oFU?7[4YU^g;6QlkZW:9nb<VOX=3[Z5n?IUfW_hnnlo6HoIT6=g>\YL\jJc[[H]Pf0lW1>NTeN8IoTREi_ClT4T5NTfcb;a945TbPk<hKRH>::^0m2;jn;hbIf9@lEXfhZo6MIF<GSGUH>bAFg2^A^Ya@XEJ8JG<U;UnldZDRmc28::di>Z@F@7>E0I2bSGi\>m4[LkBRRBR<9V=FWU6M[C9;YhB<d51g[:D0<?ni0Be`OAA^PPI7k[PER_>^6:liX=h;T[??1HCS;06foA01=2fY99o8NgKn4<in6:LDc`LJ5o>H@]TQP7Cb[QHWSngl1;KMkSK5b@8MUQE>ngcYC[2F6;N\Xa\@;IV=a9kln?QXaLFn:dRkL\dRVnB6GLnSA6\Tco25fMlC6>Rm_KD4hVZa08K4co\EW?GFh8A:K`5KU2XJRnRIgK555WMJS9SoX=LGR[KRGI^l2cCUiM_<m^G:CURKWE>=FhJi\=3iS?0Tb^;>T;Sen@JjklQg?k175^Y>F:0Z\ENI`W4__BVB9>MmT^C]99`?8iLidJjK1XiK<I@MMJD;`:Pl9=2MQMVLWmBe5Ok<hk[Pbck05h?XlYd:CI01B;VkH057d6JdO`GI>e<YC_iN?@BjC?=4X=UC\OBLA9k?=96<S]Z:;>eoEWS1ZMbkCPkFXKd826MXjiM?F]9E?JRUflfjYRZE`2B_HJ^cKKTP_WN5c0_R[ZRcQ<XiFlJoiM5A5BjBYRY_MMTnMUB>Ne=?h9`f^1o3:?E=@6do8Yj8jo@]?3oZMm]k\NE;>dC13X=hlYjBaXoG6Ul?PEiH^=:GNI3a:fbfYRSQJBa;22moP1mI_b^Uc1`DPFJDj82Zm1BI9JeI\I2nmGG2U38V;W_N;@k99DZh:Q4=ASn6DfN`kanHOLimf1IPVQG?nlWI[><b^];IG1=aLcGM^``77o@2`YZW_Rk>go84]O8^5@:LW2FOdkJGP:EA>bc;E>5AU\j221[kVA0:JQW1gQAmoMl6U\;FJ8T?i<`JY5meIIlS6A8a[K>Z3gYXc76kW`9oaQA:5hQl^<;U?^AD:Z`i^oNj;TIQURBbXMld3<ll8G?XGHm0`Qdc^4O]e4i0Nf9KfMnRjmWf9MZfMWZUUDfINa3d:g\o=VhhRS@4ETf=D8=ABn62:57DA]>D89o9JH\[kShTfWU1WWCO<ojWVN>fA<6hVCLM@6bS;MGFVRckB9R:J@De>[VcQO5a`mR3XXT9cOAY9e^^5G5RT^7F=gBL>UN`G9K\U^eb9hBE4Difcfc17iKo:6KJ0>\I3G7O<OHfdcH^99]0:X>BL7@dVl<=XeBJ0M^CIjLjC?Z;^mHcF:KI1;Ik@Za3M[=BiKAHBdSO8kNg5SMka`K9nZlnVC:mjV7iDK`kYaVM0XQK[MH7`fYlmPak:>GmKRU>7lke[B5aXhgJ24IehhhIedNJ9DJgP2a897UAf<Do6G6nn_bfU_CeIj=0EB7>Io9Z;5\J8G5fX4XR1^`oWC4eVNNkGPo5>j2N_moc4E1U9L_\;lJm53UY`U7lc;PGQS3F_W1<oolE@DUidHc:8SRdMdEEE68i<a0ABVg;L5IWhJ2iVMI3^hm7VX\X:INj<co?>jG__loQW5[:^2k[2nDWiF_YFJ[9<AnMZ2NoYFi9llG3>o;gF=mC5NE7:Z]GobUKJ`kdCdo:edekNDA[MB57B=RVSQABhg\hG1oPTiV\VJQR_:]COQFjLl_bP?bPW5E032iBgF3mL30@`Eeg05Xlma\6AdZJB=gdf<fGI0RPoPk^L@Z>WNEOHNS3AmbjeEE9A9QhE?PPPH9eNN5QCmg>CLlQQR[dFQfh7?[63:dXAijgCQRVGK5jRo_D=PNj8`d340j_Vg?E556C]<m6OFO3;@4>5V]Yn6378d<o8Uo7QdSDSmlGRSE:h<7e2@6`BbD;2<^90W8VkfHbh@Pjc`CI[8VmGNTJDHX5QYM2jdhmd2Ik7dGZIBRAZ?Rn55G^WDnAaW9gY?lJT[joUmgHmA_@lgTOcY^52OLVJoXDd;YKY@TWPKbM8[EZkG`ja@P>:lXBlXoJcbj@9fiHUh[0FE]kdYJP@ko0Q=]i6Y7Sn[SRnXKH=]dEA;eWIh1O<bR3VBU<KElYZXi;i;8O=DMCWW\g\b\3997ij_GXih1C_nYRLn`XcSZMjSW`j[?D0HdMEX9`d3AMREeNI47NBIkiKXi3DD9kfXS>2RWR7fY@QcV[B0USGO<`O@?XVn3T>_2Pbg?Mi8\HTAHkDlm9B:8bZgVJVVcO4k^mLT`BXN5OGO4j:=:SMP1:;4GJ3kIE5U3UI<=k8NLiLiGR9JcIUF\?TZa5W0JN<Rm[QZoin?TM@:3dY2FVVXh^V1Y0;TJ8EK\bAIUYb5WdWPiTTGnn@E=j0`\FGgSX`N0CfQCoR5o5>909KSZ47RJ2;G8ODg\`?o]JAi5I_V^h_h8?F=e<obGQ<\mM9[g>VEiCj_M>4A2:]gJ;KXM;1lDCi3Cb>R90[^1=7J;godT^IkB@<5LglJ2IOg:GSeITIMo=b4J>?_[ddT9EEMk_bTm2XOme=d5i;cHoi1I8h><J:Kh;DioIDASE9lDG;O??]KeD;OjZUXIPc<4?`72WaR2>L9M;ZWDlFKa=?VPkUT\[<?^M>ed9lRD1e>J;_gD[hfbTEX=VZn43ElFDF98d3f3f^ec;RLe?ZC9=XHZ0Lm2I4917eefD2LgI:i:AL\HmE2SVDE[4a\AD[[63VDEc=c4haFDE1Ii38@QH;aogo`?51ZPl4bm?fB2D[:UL9m^MG=FDW_U`JlNCQ]_;Ndk@PEmc;nf4QReIC21?[E:mFZ<PWWoZd^C24Z;2IGG^b0OlZ:POV7NZLfUY;=b<BAPEQL;OlfbX6hjnSZ@SjWX>RRHSUM4R1L7RTQo2;>KgDmJh]1H43>iTf9o5EmeBL1^[3BS[g>NVCNlbeFGJB\:d<>YG0iINGDMj;oM8J>EM1FSmE=X`:@khBVmLLg`9Zb;HnWED>[f38YloJgAPTW769U[@cB_22mU[TO^XoGc5AfS3Yno?e\Wj>D:AIF^l<CDcN`SGakK^:69;IO0W_;^[F_ef_e_mTD4bX`I[>UBDfb`?;=PjW8MnH=0]N5DRGYIj<mJUJ;F1=RUVanfB6i^fW:9VeZha4hU1m[LMdh10[>5P8JT4DBL>Wfo<9VI1MSdamA?_RQ<icF;Um=KcE[@B4VYN7]SUjoEl=E?X:Z_l1Wo>ck<KZ?MA@bB\PTLV\^XnU3E]_`Q?Gno]E>nmT78a8`CSVao_@AM@9ii:Dbfc?7oG;hoVGkB5[OJbbfk>j;53<5Ba9h6Mg0V8MHhj><_BVV=Sj:7YY@gY4;:7TEO3@hfiWM5mV[=AJc3G3NWV\Xlc8H3JC<2e??G;DiL`hS;odB1VRliQI1[5UcYiE;3R>Qf\UoSY7DFSi;=dCYNZXNZ^Mah_2U]cg=lEW`::[NdE=Gn8G;S<CM=aQQ;GoI^0aoKW?fLK@>IKZ[Feh3TdVmMFjllQM:1W1E`<;`7BojDHnIE=ogMk43`QekSJJbBGA7m<;oUG<Hn<OQniF8;PLAD1g0iV\IMBJZEFC2`ARG16IO@cC;k;Zbg_U9bEkhT87J46<aFo8od_[GFW>hOfn;=2W?`g][Q:a]S?[UJdh\I2@oVAo1AFnSh4C24O<iZ3LNZ=Ue]T;R`\ngYjLMeP1XOh9lMIU7mA@DF;RhMQ2Rm3@Ol?AdTX<Te[aYZhWm`9oZT@]VHHR:gU@<FYm<YZQ8DTYd4Md=[^a1I:^L`1CcS`7lnbCDE7Fe4QY\]n[DXcM0FL1B]M;S:;\[^If\FbK^O@78mL[[nN6Zhcf23H1l7:0K_mekeR7=9_MOnLSHd5S8GDHU`>o:PLDWDAbHQJQfDa0L7O`nhQj=KI8gT6n^j]WcQ:c:N\2jamB^J5B5A?40S>akL=[33_KNQTYUYILVE[]24Yl4d=bl0bP?m=O8g9^=h@EoHQ21>>G49G8;2_0Z[3T`7iFf;ag=Ga\8@oPB\@O\anSAjD1Tm^ZJNI1XJKbEmmFk@TSZG;T:a[cHNX8LYKkbE3RgDPRnGmiK<8oh>;eoYi_ZE<h2=OdE`dZGETY[7aF8Bff6S\2XVad84L_^>@`Dj]E\i7:535:gJ4?c@l]::]5hY6hJK_<\RkZSQoCG<Q>Vm9J>?=Xm`=io40>FcELjRL8i[2GKRj2cGWHQamhXmHQ\gE]8Yko8W_`9N4X>cfPk`mn`Q==?P]S;A7:;o:ABNlc=CE6=`e[1PMN1ZfjA3dNCPIIB7QYSDo98fdbjZn<hSo1`2`T3[DaGJ8c?_kYf:M[O4mNkW8od@bB2PSSVJ<eNDKjNaoHVo2HbSilWABVcgXa;S@>VPF[_LQFW:72F\d=nCd5[41j3X5dA8J:TVg2TL?n^KQ\ibPlf@\70cjfUS`jNcHcgiFkZ9XKAAb\^3E@PV3MaPk>JmEEnNTAGM?WkhUYnY<_F?<m:2K:Pf>X>Vb100j]dZ5g<Fd]JL[aII0ZJj@T4:io`b:<VY[MZch[h?KBmZ>Ad;@InJSYXc:7A5mY<lkTbQcmkL19QLnW[C\g5jDF@[3URHD5F@T::b4nOF\J4P`h[B:Mbi141o25P[;=_h^;3e9X>jZ:SVKO0_IK]ER4P0;Gd402UDP5>QM5Rb[K:Ug[I3cEiDXe5GKYU6`bl56UV=TcKGEMH]M`hCS?A0b<RLT<^TZ6[Z\P_UTGDC9:8RTa9CEcHG7?LFS[:^aPS;Q1@ECASi;f:>e=F7LOm`QdfMVfZLKhh2Mj0_TUoYOdR``nk0d8[K^_MGY`Q`h;ilQ;3cWP0QjHR0hRO<\6V@KHW<eghAM;m[Mf[cd]8V;5]Z1GGCDU>6R9PRYPk_V1RJWi2`V2DRKd]c[HTVW@C4^Ko5\ajBZN8Id]H\V=VPdBFJFhP0H?33CP9AlUG>aNMEYbcc3Y5O:1Re>71QmBo2];cieE4WogAj>0Ta4ID1][KZUli:`aa8jfGGN?M?VS3hC0gMJ_KS3alhd_ge9gJIUKU_CiNOCabk\OC0>Z4J;VoFhZJbSEAW0IUHJG=>LUdH;j9eRlXiKa_5on0G^VP2YYL5bJj7hHMN@OGT4biQ89<[V::ePARNCG;S`0557e0FCl51bh<382bSDn1hgOl_3dfOlZee;j;RDN3hlSgeUL4KbDgL5ei1WhE0gV>T1dSTYTXb=jJj^cMokQL=lE8WaG?C[5WJAOb8hj8Q^o]HDWUVA:o=_cg79RQeBWN[X32bNS6lEfPIQh?hJE^DFFdLKSE_VW19j2Q5>TiX=PEY7blYO7FOb32>1Y9JClF?_mUQQVd?K_T\5CGZ:@En2DAnSdLm[9Mo83FKb^HGQ0mGQjYLOCfQ@a>g9<eDh5J2mO2NFQMilVDH9RVY:=3G;I@PDK7:\<PZn:VCdE9R7GG7Ki0UG:QdP85:Z=3blDU0X4Ti>8ccLPG^a4a:Y4b2n8Pm4_UDoe9::P?gN5VIH3f6Si88;]@cI8jKWbfMRK3Q@RXQ8c3;@=UcnR7A[O9S`dm^X_7hMJSJkYm]7b:Am6\C7LHh\Qn\V3>hAEdj<g`3K35H4X8:8WE1R<;Z]@3HgEndXVQkbO`><am1FceV1:EcPlZ1J;WnknXURYFZU9Q2OFgS]38jE]ZEQ]]>oeEUo87JhTQ]h8SGVUjg;>@2lcCfDiQ1Ug<6CBBem3QJcUT19oZ]ZPe2OSZaCU>43DfC7\j40eSkcND@l<nHWHlMA`=e0ZeDijncKjCkZANX8h=UW>Zeh;=Q\J5d2Rmnb]fCkUTRiUHmj>Lhb:KI0`d=Lh9IEeGMncDKU`G:B0eT\FTegQeWU<YWPj6fhfEFWI:4aO1QC\86_aa:S51WbKURWe\^fR89a59c0a:7fBOfC?Yl_^NlTRVVMm8UG>bZO8OKkP]I:_LIlb:GA;iVo`LTj]LXfD>:mTO3l1fl@GmY@72a?Y<R405MZlV@7U1FL\jP__M[<nih0fGT;G7ZK8X]7i72:VFm^:II2:gmK5HWQ4`^TYATjV\j<JFU9l2R0^FAcCRW:h;Fg0JcUXaFU;KcP@bH?fb0VDD\PVhUeFSLHoN5_[MfWDoMZ;7g<F0WPTD1BMRZ_e_gBgP_iGjcC7Bl129?G;FDO`YfSkEMNQH?iKaW9G5^<nh\JZ@iSUA5j`cSR3ehmX[_^AWCl>CKU]22Qmma_INmVL8:2C6ZN=YY\6V5==hk\inXJ:`<b<5n4m`>5aVYVC7]FP1ZH:k_O\<GM\bXUede4@S`jgjGO59fT6aWif?:gKdZk45^85jo61j651[NkQiJ?E?BRk5fm99>X<ca[fASU9E7D73fM=W@aGoMjDSo[P6[<L;S?>2EQo_39QmgT3Na3jPFSjEN55VcH6fKlhkOhYl2m@2AQAbFkUAZ96idNYOW1cb5_>WG0?Yo?g1eikXkBKJhhd9j;\_i44iKdP8bXOK`a7S><Hfd>dXZ2V:iQRCZT>Z^Bl6_A2UMeB3\bWd`RRcEgi@d@A]8Rbb@O0I]GaBbdoBgR;8dJ]kBd238G7E7XjP@M;nc:[LN5`ocCHU<NjJgYI1`]ST:HNR\dd]gQ7fT0302B@eH9Xm73PL8dQ1oJ5f3<D?hST3Q=?3j1E1OUW4l6QU5<6J@<I@8X@Kd]Y4AfDo:SHccGR3[GY[UjK@jJIPg3P@if:Y3^;S7Qf[c7Co[L:kiYlS3@?nIUASMNE16mgKe_62mbTEVUIaF3O@4XefU@L`Zjc9eTCf6=K8]7J3RHIh;mlgf1I@9kmSLfhdS_M0<8mnW3UHKeT5Ef::3M0C^Gk:B]JP<8lQ2@=TVGRO55fn5LZ7Dh5HMPCoBWO"), []byte{}, 256)
	//	fmt.Println(string(decryptBytes))
	//	fmt.Println(err)

	//	gotrixHandler = handler.SimpleHandler{}
	//	gotrixHandler.Init()
	//
	//	checkedParams := &global.CheckedParams{Func: 2003, V: make(map[string]interface{})}
	//	checkedParams.V["userid"] = float64(2)
	//	response, err := gotrixHandler.Handle(checkedParams)
	//	log.Println(response)
	//	log.Println(err)

	//	for i, args := range os.Args {
	//		switch args {
	//		case "--decrypt", "-d":
	//			global.Config.Args.Decrypt = true
	//			break
	//		case "--console", "-c":
	//			global.Config.Args.Console = true
	//			break
	//		case "--password", "-p":
	//			global.Config.Args.Password = os.Args[i+1]
	//			break
	//		}
	//	}
	//
	//	if len(global.Config.Args.Password) == 0 {
	//		global.InitPassword()
	//	}
	//
	//	global.InitConfiguration()
	//
	//	for _, args := range os.Args {
	//		switch args {
	//		case "start":
	//			if global.Config.Args.Console {
	//				GotrixServer()
	//			} else {
	//				filePath, _ := filepath.Abs(os.Args[0])
	//				args := append(os.Args, "--console", "--password", global.Config.Args.Password)
	//				logFile, _ := os.Create(global.Config.LogFile)
	//				process, err := os.StartProcess(filePath, args, &os.ProcAttr{Files: []*os.File{logFile, logFile, logFile}})
	//				if err != nil {
	//					log.Println(err)
	//				}
	//				log.Println(process)
	//			}
	//			break
	//		}
	//	}

}

var gotrixChecker global.Checker
var gotrixHandler global.Handler

func GotrixServer() {

	// -----杀掉原有实例---------------------------------------------------------
	c := exec.Command("netstat", "/ano")
	bs, err := c.Output()
	if err != nil {
		fmt.Println(err)
	}

	reg := regexp.MustCompile("TCP\\s*0\\.0\\.0\\.0:9080.*LISTENING\\s*(\\d*)")
	matches := reg.FindSubmatch(bs)

	if matches != nil && len(matches) > 0 {
		c1 := exec.Command("taskkill", "-f", "/pid", string(matches[1]))
		err1 := c1.Run()
		if err1 != nil {
			fmt.Println(err1)
		}
	}
	// -----------------------------------------------------------------------

	gotrixChecker = checker.EncryptChecker{}
	gotrixHandler = handler.SimpleHandler{}
	gotrixHandler.Init()

	http.HandleFunc("/gotrix/", serverHandler)
	http.HandleFunc("/gotrix/wxpay.action", wxpayCallback)
	http.Handle("/", http.FileServer(http.Dir("src/github.com/zhutingle/gotrix/static")))

	err = http.ListenAndServe(":9080", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func writeError(w http.ResponseWriter, err *global.GotrixError) {
	if err.Status > 0 {
		w.Write([]byte(fmt.Sprintf("{\"status\":%d,\"msg\":\"%s\"}", err.Status, err.Msg)))
	} else {
		w.Write([]byte(err.Msg))
	}
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	var start = time.Now().UnixNano()
	var logBuffer bytes.Buffer
	// --------------------参数解析器--------------------
	checkedParams, gErr := gotrixChecker.Check(r, gotrixHandler)
	if gErr != nil {
		writeError(w, gErr)

		logBuffer.WriteString("\n----Error: ")
		logBuffer.WriteString(fmt.Sprint(gErr))
		logBuffer.WriteRune('\n')
		log.Println(logBuffer.String())
		return
	}

	logBuffer.WriteString("\n----Func: ")
	logBuffer.WriteString(strconv.FormatInt(int64(checkedParams.Func), 10))
	logBuffer.WriteString("\n----Param: ")
	logBuffer.WriteString(fmt.Sprint(checkedParams.V))

	// --------------------业务执行器--------------------
	var response interface{}
	response, gErr = gotrixHandler.Handle(checkedParams)
	if gErr != nil {
		writeError(w, gErr)

		logBuffer.WriteString("\n----Error: ")
		logBuffer.WriteString(fmt.Sprint(gErr))
		logBuffer.WriteRune('\n')
		log.Println(logBuffer.String())
		return
	}

	// --------------------结果输出器--------------------
	buffer := bytes.NewBufferString("{\"status\":0,\"msg\":\"成功\",\"data\":")
	str, _ := json.Marshal(response)
	buffer.Write(str)
	buffer.WriteString("}")
	encryptResult, e := global.AesEncrypt(buffer.Bytes(), checkedParams.Pass, 256)
	if e != nil {
		writeError(w, global.RETURN_DATE_ECNRYPT_ERROR)

		logBuffer.WriteString("\n----Error: ")
		logBuffer.WriteString(fmt.Sprint(e))
		logBuffer.WriteRune('\n')
		log.Println(logBuffer.String())
		return
	}
	w.Write(encryptResult)

	logBuffer.WriteString("\n----Result: ")
	logBuffer.Write(str)
	logBuffer.WriteString("\n----Spend: ")
	logBuffer.WriteString(strconv.FormatInt((time.Now().UnixNano()-start)/1000000, 10))
	logBuffer.WriteString(" ms")
	logBuffer.WriteRune('\n')
	log.Println(logBuffer.String())

}

func wxpayCallback(w http.ResponseWriter, r *http.Request) {
	checkedParams := &global.CheckedParams{Func: 1001, V: make(map[string]interface{}, 0)}
	weichat, err := weichat.WxpayCallback(w, r)
	if err != nil {
		fmt.Println(err)
	}
	checkedParams.V["weichat"] = weichat
	fmt.Println(checkedParams.V)
	response, gErr := gotrixHandler.Handle(checkedParams)
	if gErr != nil {
		writeError(w, gErr)
	} else {
		writeError(w, &global.GotrixError{Status: 0, Msg: fmt.Sprintf("%v", response)})
	}
}
