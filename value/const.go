// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains 3011-digit string representations of fundamental constants.
// They are stored in 10000-bit Floats, almost sufficient to represent the values.
// When used in calculations, those Floats are truncated to the precision of the operands.

// In the unlikely event more digits are needed, it's easy to find the values online and
// lengthen the constants. Or compute them: By a magic of math, log 2 is not needed
// to compute log 2 by the algorithm here (the exponent is zero in the Taylor-Maclaurin
// series) so it is possible to bootstrap to huge precisions.

package value

import (
	"fmt"
	"math/big"

	"robpike.io/ivy/config"
)

const (
	// 10000*log(2)/log(10) == 3011
	constPrecisionInBits   = 10000
	constPrecisionInDigits = 3011
)

var (
	// set to constPrecision
	floatE          *big.Float
	floatPi         *big.Float
	floatPiBy2      *big.Float
	floatMinusPiBy2 *big.Float
	floatLog2       *big.Float
	floatLog10      *big.Float
	floatZero       *big.Float
	floatOne        *big.Float
	floatTwo        *big.Float
	floatHalf       *big.Float
	floatMinusOne   *big.Float

	bigIntTen     = big.NewInt(10)
	bigIntBillion = big.NewInt(1e9)
	MaxBigInt63   = big.NewInt(int64(^uint64(0) >> 1)) // Used in ../parse/special.go

	bigRatOne      = big.NewRat(1, 1)
	bigRatMinusOne = big.NewRat(-1, 1)
	bigRatTen      = big.NewRat(10, 1)
	bigRatBillion  = big.NewRat(1e9, 1)
)

const strE = "2.7182818284590452353602874713526624977572470936999595749669676277240766303535475945713821785251664274274663919320030599218174135966290435729003342952605956307381323286279434907632338298807531952510190115738341879307021540891499348841675092447614606680822648001684774118537423454424371075390777449920695517027618386062613313845830007520449338265602976067371132007093287091274437470472306969772093101416928368190255151086574637721112523897844250569536967707854499699679468644549059879316368892300987931277361782154249992295763514822082698951936680331825288693984964651058209392398294887933203625094431173012381970684161403970198376793206832823764648042953118023287825098194558153017567173613320698112509961818815930416903515988885193458072738667385894228792284998920868058257492796104841984443634632449684875602336248270419786232090021609902353043699418491463140934317381436405462531520961836908887070167683964243781405927145635490613031072085103837505101157477041718986106873969655212671546889570350354021234078498193343210681701210056278802351930332247450158539047304199577770935036604169973297250886876966403555707162268447162560798826517871341951246652010305921236677194325278675398558944896970964097545918569563802363701621120477427228364896134225164450781824423529486363721417402388934412479635743702637552944483379980161254922785092577825620926226483262779333865664816277251640191059004916449982893150566047258027786318641551956532442586982946959308019152987211725563475463964479101459040905862984967912874068705048958586717479854667757573205681288459205413340539220001137863009455606881667400169842055804033637953764520304024322566135278369511778838638744396625322498506549958862342818997077332761717839280349465014345588970719425863987727547109629537415211151368350627526023264847287039207643100595841166120545297030236472549296669381151373227536450988890313602057248176585118063036442812314965507047510254465011727211555194866850800368532281831521960037356252794495158284188294787610852639813955990067376482922443752871846245780361929819713991475644882626039033814418232625150974827987779964373089970388867782271383605772978824125611907176639465070633045279546618550966661856647097113444740160704626215680717481877844371436988218559670959102596862002353718588748569652200050311734392073211390803293634479727355955277349071783793421637012050054513263835440001863239914907054797780566978533580489669062951194324730995876552368128590413832411607226029983305353708761389396391779574540161372236187893652605381558415871869255386061647798340254351284396129460352913325942794904337299085731580290958631382683291477116396337092400316894586360606458459251269946557248391865642097526850823075442545993769170419777800853627309417101634349076964237222943523661255725088147792231519747780605696725380171807763603462459278778465850656050780844211529697521890874019660906651803516501792504619501366585436632712549639908549144200014574760819302212066024330096412704894390397177195180699086998606636583232278709376502260"

const strPi = "3.1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679821480865132823066470938446095505822317253594081284811174502841027019385211055596446229489549303819644288109756659334461284756482337867831652712019091456485669234603486104543266482133936072602491412737245870066063155881748815209209628292540917153643678925903600113305305488204665213841469519415116094330572703657595919530921861173819326117931051185480744623799627495673518857527248912279381830119491298336733624406566430860213949463952247371907021798609437027705392171762931767523846748184676694051320005681271452635608277857713427577896091736371787214684409012249534301465495853710507922796892589235420199561121290219608640344181598136297747713099605187072113499999983729780499510597317328160963185950244594553469083026425223082533446850352619311881710100031378387528865875332083814206171776691473035982534904287554687311595628638823537875937519577818577805321712268066130019278766111959092164201989380952572010654858632788659361533818279682303019520353018529689957736225994138912497217752834791315155748572424541506959508295331168617278558890750983817546374649393192550604009277016711390098488240128583616035637076601047101819429555961989467678374494482553797747268471040475346462080466842590694912933136770289891521047521620569660240580381501935112533824300355876402474964732639141992726042699227967823547816360093417216412199245863150302861829745557067498385054945885869269956909272107975093029553211653449872027559602364806654991198818347977535663698074265425278625518184175746728909777727938000816470600161452491921732172147723501414419735685481613611573525521334757418494684385233239073941433345477624168625189835694855620992192221842725502542568876717904946016534668049886272327917860857843838279679766814541009538837863609506800642251252051173929848960841284886269456042419652850222106611863067442786220391949450471237137869609563643719172874677646575739624138908658326459958133904780275900994657640789512694683983525957098258226205224894077267194782684826014769909026401363944374553050682034962524517493996514314298091906592509372216964615157098583874105978859597729754989301617539284681382686838689427741559918559252459539594310499725246808459872736446958486538367362226260991246080512438843904512441365497627807977156914359977001296160894416948685558484063534220722258284886481584560285060168427394522674676788952521385225499546667278239864565961163548862305774564980355936345681743241125150760694794510965960940252288797108931456691368672287489405601015033086179286809208747609178249385890097149096759852613655497818931297848216829989487226588048575640142704775551323796414515237462343645428584447952658678210511413547357395231134271661021359695362314429524849371871101457654035902799344037420073105785390621983874478084784896833214457138687519435064302184531910484810053706146806749192781911979399520614196634287544406437451237181921799983910159195618146751426912397489409071864942319615679452080"

const strLog2 = "0.6931471805599453094172321214581765680755001343602552541206800094933936219696947156058633269964186875420014810205706857336855202357581305570326707516350759619307275708283714351903070386238916734711233501153644979552391204751726815749320651555247341395258829504530070953263666426541042391578149520437404303855008019441706416715186447128399681717845469570262716310645461502572074024816377733896385506952606683411372738737229289564935470257626520988596932019650585547647033067936544325476327449512504060694381471046899465062201677204245245296126879465461931651746813926725041038025462596568691441928716082938031727143677826548775664850856740776484514644399404614226031930967354025744460703080960850474866385231381816767514386674766478908814371419854942315199735488037516586127535291661000710535582498794147295092931138971559982056543928717000721808576102523688921324497138932037843935308877482597017155910708823683627589842589185353024363421436706118923678919237231467232172053401649256872747782344535347648114941864238677677440606956265737960086707625719918473402265146283790488306203306114463007371948900274364396500258093651944304119115060809487930678651588709006052034684297361938412896525565396860221941229242075743217574890977067526871158170511370091589426654785959648906530584602586683829400228330053820740056770530467870018416240441883323279838634900156312188956065055315127219939833203075140842609147900126516824344389357247278820548627155274187724300248979454019618723398086083166481149093066751933931289043164137068139777649817697486890388778999129650361927071088926410523092478391737350122984242049956893599220660220465494151061391878857442455775102068370308666194808964121868077902081815885800016881159730561866761991873952007667192145922367206025395954365416553112951759899400560003665135675690512459268257439464831683326249018038242408242314523061409638057007025513877026817851630690255137032340538021450190153740295099422629957796474271381573638017298739407042421799722669629799393127069357472404933865308797587216996451294464918837711567016785988049818388967841349383140140731664727653276359192335112333893387095132090592721854713289754707978913844454666761927028855334234298993218037691549733402675467588732367783429161918104301160916952655478597328917635455567428638774639871019124317542558883012067792102803412068797591430812833072303008834947057924965910058600123415617574132724659430684354652111350215443415399553818565227502214245664400062761833032064727257219751529082785684213207959886389672771195522188190466039570097747065126195052789322960889314056254334425523920620303439417773579455921259019925591148440242390125542590031295370519220615064345837878730020354144217857580132364516607099143831450049858966885772221486528821694181270488607589722032166631283783291567630749872985746389282693735098407780493950049339987626475507031622161390348452994249172483734061366226383493681116841670569252147513839306384553718626877973288955588716344297562447553923663694888778238901749810273565524050"

const strLog10 = "2.3025850929940456840179914546843642076011014886287729760333279009675726096773524802359972050895982983419677840422862486334095254650828067566662873690987816894829072083255546808437998948262331985283935053089653777326288461633662222876982198867465436674744042432743651550489343149393914796194044002221051017141748003688084012647080685567743216228355220114804663715659121373450747856947683463616792101806445070648000277502684916746550586856935673420670581136429224554405758925724208241314695689016758940256776311356919292033376587141660230105703089634572075440370847469940168269282808481184289314848524948644871927809676271275775397027668605952496716674183485704422507197965004714951050492214776567636938662976979522110718264549734772662425709429322582798502585509785265383207606726317164309505995087807523710333101197857547331541421808427543863591778117054309827482385045648019095610299291824318237525357709750539565187697510374970888692180205189339507238539205144634197265287286965110862571492198849978748873771345686209167058498078280597511938544450099781311469159346662410718466923101075984383191912922307925037472986509290098803919417026544168163357275557031515961135648465461908970428197633658369837163289821744073660091621778505417792763677311450417821376601110107310423978325218948988175979217986663943195239368559164471182467532456309125287783309636042629821530408745609277607266413547875766162629265682987049579549139549180492090694385807900327630179415031178668620924085379498612649334793548717374516758095370882810674524401058924449764796860751202757241818749893959716431055188481952883307466993178146349300003212003277656541304726218839705967944579434683432183953044148448037013057536742621536755798147704580314136377932362915601281853364984669422614652064599420729171193706024449293580370077189810973625332245483669885055282859661928050984471751985036666808749704969822732202448233430971691111368135884186965493237149969419796878030088504089796185987565798948364452120436982164152929878117429733325886079159125109671875109292484750239305726654462762009230687915181358034777012955936462984123664970233551745861955647724618577173693684046765770478743197805738532718109338834963388130699455693993461010907456160333122479493604553618491233330637047517248712763791409243983318101647378233796922656376820717069358463945316169494117018419381194054164494661112747128197058177832938417422314099300229115023621921867233372683856882735333719251034129307056325444266114297653883018223840910261985828884335874559604530045483707890525784731662837019533922310475275649981192287427897137157132283196410034221242100821806795252766898581809561192083917607210809199234615169525990994737827806481280587927319938934534153201859697110214075422827962982370689417647406422257572124553925261793736524344405605953365915391603125244801493132345724538795243890368392364505078817313597112381453237015084134911223243909276817247496079557991513639828810582857405380006533716555530141963322419180876210182049194926514838926922937079"

func newF(conf *config.Config) *big.Float {
	return new(big.Float).SetPrec(conf.FloatPrec())
}

func newFloat(c Context) *big.Float {
	return newF(c.Config())
}

func Consts(c Context) (e, pi BigFloat) {
	conf := c.Config()
	if conf.FloatPrec() > constPrecisionInBits {
		fmt.Fprintf(c.Config().ErrOutput(), "warning: precision too high; only have %d digits (%d bits) of precision for e and pi", constPrecisionInDigits, constPrecisionInBits)
	}
	floatZero = newF(conf).SetInt64(0)
	floatOne = newF(conf).SetInt64(1)
	floatTwo = newF(conf).SetInt64(2)
	floatHalf = newF(conf).SetFloat64(0.5)
	floatMinusOne = newF(conf).SetInt64(-1)
	var ok bool
	floatE, ok = newF(conf).SetString(strE)
	if !ok {
		panic("setting e")
	}
	floatPi, ok = newF(conf).SetString(strPi)
	if !ok {
		panic("setting pi")
	}
	floatPiBy2, ok = newF(conf).SetString(strPi)
	if !ok {
		panic("setting pi")
	}
	floatPiBy2.Quo(floatPiBy2, floatTwo)
	floatMinusPiBy2 = newF(conf).Neg(floatPiBy2)
	floatLog2, ok = newF(conf).SetString(strLog2)
	if !ok {
		panic("setting log(2)")
	}
	floatLog10, ok = newF(conf).SetString(strLog10)
	if !ok {
		panic("setting log(10)")
	}
	return BigFloat{newF(conf).Set(floatE)}, BigFloat{newF(conf).Set(floatPi)}

}
